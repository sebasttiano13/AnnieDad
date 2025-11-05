package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sebasttiano13/AnnieDad/internal/repository"
	"github.com/sebasttiano13/AnnieDad/pkg/jwt"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

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
		"/main.AuthService/LoginTelegram":    true,
		"/main.AuthService/RefreshToken":     true,
		"/main.AuthService/RegisterTelegram": true,
		"/main.AuthService/LinkWeb":          true,
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
	claims, err := i.jwtManager.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	// Set user id to metadata
	md, ok = metadata.FromIncomingContext(ctx)
	if ok {
		md.Append("user-id", claims.ID)
	}

	newCtx := metadata.NewIncomingContext(ctx, md)
	return newCtx, nil
}

type ApiKeyChecker interface {
	Validate(ctx context.Context, token string) (bool, error)
}
type ApiKeyInterceptor struct {
	checker ApiKeyChecker
}

func NewApiKeyInterceptor(checker ApiKeyChecker) *ApiKeyInterceptor {
	return &ApiKeyInterceptor{checker: checker}
}

func (i *ApiKeyInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		apiKey, err := getApiClientTokenFromContext(ctx)
		if err != nil {
			if errors.Is(err, ErrGetApiTokenFromContext) {
				// No api key meta
				logger.Errorf(err.Error())
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
			logger.Debugf("api key check error: %v", err)
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		valid, err := i.checker.Validate(ctx, apiKey)
		if err != nil {
			if errors.Is(err, repository.ErrBadAPIKeyFormat) {
				logger.Debugf("api key %s check error: %v", apiKey, err)
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			logger.Debugf("api key check error: %v", err)
			return nil, status.Error(codes.Internal, "failed to check api key")
		}
		if !valid {
			logger.Debugf("api key %s is invalid", apiKey)
			return nil, status.Error(codes.PermissionDenied, "invalid api key")
		}

		newCtx := context.WithValue(ctx, "api-client", true)
		return handler(newCtx, req)
	}
}
