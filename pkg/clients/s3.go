package clients

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	ErrS3Client             = errors.New("s3 client failed")
	ErrS3ClientGetUploadURL = errors.New("failed to get presign upload URL")
)

// S3Client simple client
type S3Client struct {
	client        *s3.Client
	presigner     *s3.PresignClient
	ExpireURLTime time.Duration
}

// NewS3Client constructor for S3Client
func NewS3Client(cfg aws.Config, expires time.Duration) *S3Client {
	client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)
	return &S3Client{client: client, ExpireURLTime: expires, presigner: presigner}
}

func (s *S3Client) GetUploadURL(ctx context.Context, bucket string, key string) (string, error) {
	presignedReq, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(s.ExpireURLTime))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrS3ClientGetUploadURL, err)
	}

	return presignedReq.URL, nil
}
