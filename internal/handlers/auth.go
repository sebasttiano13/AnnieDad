package handlers

import (
	"context"
	"errors"

	pb "github.com/sebasttiano13/AnnieDad/internal/proto/auth"
	"github.com/sebasttiano13/AnnieDad/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServer) LoginBot(ctx context.Context, in *pb.BotLoginRequest) (*pb.AuthResponse, error) {
	accessToken, refreshToken, err := s.BotAuth.LoginBot(ctx, in.TelegramId)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, `user with telegram id %d not found`, in.GetTelegramId())
		}
		return nil, status.Error(codes.Internal, `login failed`)
	}
	return &pb.AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, in *pb.RefreshRequest) (*pb.AuthResponse, error) {
	refreshToken := in.GetRefreshToken()
	accessToken, refreshToken, err := s.Refresh.Refresh(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInternalAuthService) {
			return nil, status.Error(codes.Internal, "internal server error")
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	out := &pb.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return out, nil
}
func (s *AuthServer) LinkWeb(ctx context.Context, in *pb.LinkWebRequest) (*emptypb.Empty, error) {
	return nil, nil
}
