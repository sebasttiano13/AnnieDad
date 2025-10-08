package handlers

import (
	"context"

	pb "github.com/sebasttiano13/AnnieDad/internal/proto"
)

func (m *MediaServer) PostURL(ctx context.Context, in *pb.PostMediaRequest) (*pb.PostMediaResponse, error) {

	url, err := m.Media.PostURL(ctx, "my_file")
	if err != nil {
		return nil, err
	}
	return &pb.PostMediaResponse{Url: url}, nil
}

func (m *MediaServer) GetURL(ctx context.Context, in *pb.GetMediaRequest) (*pb.GetMediaResponse, error) {
	url, err := m.Media.GetUploadURL(ctx, in.Filename)
	if err != nil {
		return nil, err
	}
	return &pb.GetMediaResponse{Url: url}, nil
}
