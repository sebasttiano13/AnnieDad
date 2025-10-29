package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sebasttiano13/AnnieDad/internal/models"
	"github.com/sebasttiano13/AnnieDad/pkg/jwt"
)

// pgError алиас для *pgconn.PgError
var pgError *pgconn.PgError

var ErrInternalAuthService = errors.New("auth service internal error")

type AuthRepo interface {
	GetByUsername(ctx context.Context, user *models.User) error
	GetByTelegramID(ctx context.Context, user *models.User) error
	AddUser(ctx context.Context, user *models.User) error
	AddBotUser(ctx context.Context, user *models.User) error
	LinkTelegramUser(ctx context.Context, user *models.User) error
}

type MediaRepo interface{}

type AuthService struct {
	Repo       AuthRepo
	jwtManager *jwt.JWTManager
}

func NewAuthService(repo AuthRepo, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{Repo: repo, jwtManager: jwtManager}
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
