package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sebasttiano13/AnnieDad/internal/domains"
	"github.com/sebasttiano13/AnnieDad/pkg/jwt"
)

// pgError алиас для *pgconn.PgError
var pgError *pgconn.PgError

var ErrInternalAuthService = errors.New("auth service internal error")

type TokensOpts struct {
	AccessTokenDuration     time.Duration
	RefreshTokenDuration    time.Duration
	RefreshTokenRenewBefore time.Duration
	RefreshCleanupInterval  time.Duration
}

type AuthRepo interface {
	GetByUsername(ctx context.Context, userName string) (*domains.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*domains.User, error)
	AddUser(ctx context.Context, userName string, userPassword string) (*domains.User, error)
	AddBotUser(ctx context.Context, telegramID int64, userName string) (*domains.User, error)
	LinkTelegramUser(ctx context.Context, user *domains.User) error
}

type RefreshRepo interface {
	SaveRefresh(ctx context.Context, userID, token, tokenID string, expiresAt time.Time) error
	GetRefresh(ctx context.Context, id string) (*domains.RefreshToken, error)
	DeleteRefresh(ctx context.Context, id string) error
	DeleteExpiredRefresh(ctx context.Context) (int64, error)
}

type MediaRepo interface{}

type TokenService struct {
	repo       RefreshRepo
	jwtManager *jwt.JWTManager
	opts       *TokensOpts
}

func NewTokenService(repo RefreshRepo, jwtManager *jwt.JWTManager, opts *TokensOpts) *TokenService {
	return &TokenService{repo: repo, jwtManager: jwtManager, opts: opts}
}

type AuthService struct {
	Repo   AuthRepo
	Tokens TokenService
}

func NewAuthService(repo AuthRepo, TS *TokenService) *AuthService {
	return &AuthService{Repo: repo, Tokens: *TS}
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
