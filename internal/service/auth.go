package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/sebasttiano13/AnnieDad/internal/repository"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrLoginFailed            = errors.New("login failed")
	ErrUserRegistrationFailed = errors.New("user registration failed")
	ErrUserAlreadyExist       = errors.New("user already exist")
)

// RegisterWeb register new user via web
func (a *AuthService) RegisterWeb(ctx context.Context, name, password string) (string, string, error) {
	return "", "", nil
}

// RegisterTelegram register new user via telegram id
func (a *AuthService) RegisterTelegram(ctx context.Context, telegramID int64, name string) (string, string, error) {
	user, err := a.Repo.AddBotUser(ctx, telegramID, name)
	if err != nil {
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				logger.Debugf("user with telegram ID %d already exists", user.TelegramID)
				return "", "", ErrUserAlreadyExist
			}
		}
		logger.Errorf("failed to add user with telegram ID %d: %v", user.TelegramID, err)
		return "", "", ErrUserRegistrationFailed
	}
	logger.Infof("%s registered successfully", user.Name)
	access, refresh, err := a.Tokens.GetTokens(ctx, user.ID)
	if err != nil {
		logger.Errorf("failed to get tokens: %v", err)
		return "", "", ErrUserRegistrationFailed
	}
	return access, refresh, nil
}

func (a *AuthService) LoginWeb(ctx context.Context, name, password string) (string, string, error) {
	return "", "", nil
}

func (a *AuthService) LinkTelegramBot(ctx context.Context, bindToken string) error {
	return nil
}

func (a *AuthService) LoginTelegram(ctx context.Context, telegramID int64) (string, string, error) {
	user, err := a.Repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if errors.Is(err, repository.ErrDBNoRows) {
			logger.Infof("user by telegram id %d not found", telegramID)
			return "", "", ErrUserNotFound
		}
		logger.Errorf("user with telegram id %v login failed: %v", telegramID, err)
		return "", "", fmt.Errorf("%w: %v", ErrLoginFailed, err)
	}

	return a.Tokens.GetTokens(ctx, user.ID)
}

func (a *AuthService) LinkWeb(ctx context.Context, bindToken string) error {
	return nil
}
func (a *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	return a.Tokens.RefreshTokens(ctx, refreshToken)
}
