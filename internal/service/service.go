package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sebasttiano13/AnnieDad/internal/models"
)

// pgError алиас для *pgconn.PgError
var pgError *pgconn.PgError

type AuthRepo interface {
	GetUser(ctx context.Context, user *models.User) error
	AddUser(ctx context.Context, user *models.User) error
}

type MediaRepo interface{}

type AuthService struct {
	Repo AuthRepo
}

func NewAuthService(repo AuthRepo) *AuthService {
	return &AuthService{Repo: repo}
}

type MediaService struct {
	Repo MediaRepo
	S3   S3
}

func NewMediaService(repo MediaRepo, s3client S3) *MediaService {
	return &MediaService{Repo: repo, S3: s3client}
}

type S3 interface {
	DownloadURL(ctx context.Context, bucket string, key string) (string, error)
	UploadURL(ctx context.Context, bucket string, key string) (string, error)
}
