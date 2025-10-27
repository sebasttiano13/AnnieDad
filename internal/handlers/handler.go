package handlers

import (
	"context"
	"errors"

	pbMedia "github.com/sebasttiano13/AnnieDad/internal/proto/anniedad"
	pbAuth "github.com/sebasttiano13/AnnieDad/internal/proto/auth"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrInternalGrpc       = errors.New("internal grpc server error")
	ErrGetUserFromContext = errors.New("failed get user id from context")
)

type AuthServer struct {
	BotAuth BotAuthenticator
	Refresh Refresher
	pbAuth.UnimplementedAuthServiceServer
}

type MediaServer struct {
	Media MediaServ
	pbMedia.UnimplementedMediaServer
}

type WebAuthenticator interface {
	Register(ctx context.Context, name, password string) (string, string, error)
	LoginWeb(ctx context.Context, name, password string) (string, string, error)
	LinkTelegramBot(ctx context.Context, bindToken string) error
}

type BotAuthenticator interface {
	LoginBot(ctx context.Context, telegramID int64) (string, string, error)
	LinkWeb(ctx context.Context, bindToken string) error
}

type Refresher interface {
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}

type MediaServ interface {
	PostURL(ctx context.Context, fileName string) (string, error)
	GetUploadURL(ctx context.Context, fileName string) (string, error)
}

var getUserIDFromContext = func(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error(ErrGetUserFromContext.Error())
		return "", status.Errorf(codes.Internal, ErrInternalGrpc.Error())
	}
	userID := md.Get("user-id")[0]
	return userID, nil
}
