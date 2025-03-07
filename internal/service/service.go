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

type AuthService struct {
	Repo AuthRepo
}

func NewAuthService(repo AuthRepo) *AuthService {
	return &AuthService{Repo: repo}
}

func (a *AuthService) Register(ctx context.Context, name, password string) error {
	return nil
}

func (a *AuthService) Login(ctx context.Context, name, password string) (int, error) {
	return 0, nil
}
