package handlers

import (
	"context"
	"errors"

	pbMedia "github.com/sebasttiano13/AnnieDad/internal/proto/anniedad"
	pbAuth "github.com/sebasttiano13/AnnieDad/internal/proto/auth"
	"google.golang.org/grpc/metadata"
)

var (
	ErrInternalGrpc           = errors.New("internal grpc server error")
	ErrGetUserFromContext     = errors.New("failed get user id from context")
	ErrGetApiTokenFromContext = errors.New("failed get api client token from context")
	ErrNoMetadata             = errors.New("metadata is not provided")
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
	RegisterWeb(ctx context.Context, name, password string) (string, string, error)
	LoginWeb(ctx context.Context, name, password string) (string, string, error)
	LinkTelegramBot(ctx context.Context, bindToken string) error
}

type BotAuthenticator interface {
	RegisterTelegram(ctx context.Context, telegramID int64, name string) (string, string, error)
	LoginTelegram(ctx context.Context, telegramID int64) (string, string, error)
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
		return "", ErrNoMetadata
	}
	keys := md.Get("keys")
	if len(keys) == 0 {
		return "", ErrGetUserFromContext
	}
	userID := keys[0]
	return userID, nil
}

var getApiClientTokenFromContext = func(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrNoMetadata
	}
	key := md.Get("x-api-key")
	if len(key) == 0 {
		return "", ErrGetApiTokenFromContext
	}
	token := key[0]
	return token, nil
}
