package service

import (
	"context"
	"fmt"
)

func (m *MediaService) PostURL(ctx context.Context, fileName string) (string, error) {
	url, err := m.S3.GetUploadURL(ctx, "annie", "avatar.png")
	if err != nil {
		fmt.Println(err)
		return "example.com", nil
	}
	return url, nil
}
