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
}

func NewMediaService(repo MediaRepo) *MediaService {
	return &MediaService{Repo: repo}
}
