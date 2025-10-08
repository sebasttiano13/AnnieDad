package service

import (
	"context"

	"github.com/sebasttiano13/AnnieDad/pkg/logger"
)

func (m *MediaService) PostURL(ctx context.Context, fileName string) (string, error) {
	url, err := m.S3.UploadURL(ctx, "annie", fileName)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	return url, nil
}

func (m *MediaService) GetUploadURL(ctx context.Context, fileName string) (string, error) {
	url, err := m.S3.DownloadURL(ctx, "annie", fileName)
	if err != nil {
		logger.Error(err.Error())
		return "example.com", nil
	}
	return url, nil
}
