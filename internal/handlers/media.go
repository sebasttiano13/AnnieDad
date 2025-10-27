package handlers

import (
	"context"
	"fmt"

	pb "github.com/sebasttiano13/AnnieDad/internal/proto/anniedad"
)

func (m *MediaServer) PostURL(ctx context.Context, in *pb.PostMediaRequest) (*pb.PostMediaResponse, error) {

	userID, err := getUserIDFromContext(ctx)
	fmt.Println("userID:", userID)

	url, err := m.Media.PostURL(ctx, in.Filename)
	if err != nil {
		return nil, err
	}
	return &pb.PostMediaResponse{Url: url}, nil
}

func (m *MediaServer) GetURL(ctx context.Context, in *pb.GetMediaRequest) (*pb.GetMediaResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	fmt.Println("userID:", userID)

	url, err := m.Media.GetUploadURL(ctx, in.Filename)
	if err != nil {
		return nil, err
	}
	return &pb.GetMediaResponse{Url: url}, nil
}
