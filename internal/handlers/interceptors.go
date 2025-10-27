package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sebasttiano13/AnnieDad/pkg/jwt"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var ErrNoMetadata = errors.New("metadata is not provided")
var ErrNoAccessToken = errors.New("authorization token is not provided")

// InterceptorLogger logging all incoming requests
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		attrs := make([]any, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key, ok := fields[i].(string)
			if !ok || i+1 >= len(fields) {
				continue
			}

			attrs = append(attrs, slog.Any(key, fields[i+1]))
		}

		logger := l.With(attrs...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			logger.Info(msg)
		}
	})
}

type AuthInterceptor struct {
	jwtManager       *jwt.JWTManager
	whitelistMethods map[string]bool
}

func NewAuthInterceptor(jwtManager *jwt.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, map[string]bool{
		"/main.Auth/LoginBot":     true,
		"/main.Auth/RefreshToken": true,
		"/main.Auth/LinkWeb":      true,
	}}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Infof("--> unary interceptor: %s", info.FullMethod)
		newCtx, err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		if newCtx != nil {
			ctx = newCtx
		}
		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	_, ok := i.whitelistMethods[method]
	if ok {
		// everyone can access
		return nil, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, ErrNoMetadata.Error())
	}
	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, ErrNoAccessToken.Error())
	}

	accessToken := values[0]
	claims, err := i.jwtManager.VerifyToken(accessToken, jwt.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%s: %v", jwt.ErrInvalidAccessToken, err)
	}

	// Set user id to metadata
	md, ok = metadata.FromIncomingContext(ctx)
	if ok {
		md.Append("user-id", claims.ID)
	}

	newCtx := metadata.NewIncomingContext(ctx, md)
	return newCtx, nil
}
