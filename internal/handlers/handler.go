package handlers

import (
	"context"
	"errors"

	pb "github.com/sebasttiano13/AnnieDad/internal/proto"
)

var (
	ErrInternalGrpc = errors.New("internal grpc server error")
)

type AuthServer struct {
	Auth Authenticator
	pb.UnimplementedAuthServer
}

type MediaServer struct {
	Media MediaServ
	pb.UnimplementedPresignedURLServer
}

type Authenticator interface {
	Register(ctx context.Context, name, password string) error
	Login(ctx context.Context, name, password string) (int, error)
}

type MediaServ interface {
	PostURL(ctx context.Context, fileName string) (string, error)
}
