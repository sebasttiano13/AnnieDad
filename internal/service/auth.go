package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/sebasttiano13/AnnieDad/internal/models"
	"github.com/sebasttiano13/AnnieDad/internal/repository"
	"github.com/sebasttiano13/AnnieDad/pkg/jwt"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrLoginFailed           = errors.New("login failed")
	ErrRefreshTokenFailed    = errors.New("refresh login failed")
	ErrUserRegisrationFailed = errors.New("user registration failed")
	ErrUserAlreadyExist      = errors.New("user already exist")
)

// RegisterWeb register new user via web
func (a *AuthService) RegisterWeb(ctx context.Context, name, password string) (string, string, error) {
	return "", "", nil
}

// RegisterTelegram register new user via telegram id
func (a *AuthService) RegisterTelegram(ctx context.Context, telegramID int64, name string) (string, string, error) {
	user := &models.User{TelegramID: telegramID, Name: sql.NullString{String: name, Valid: true}}
	if err := a.Repo.AddBotUser(ctx, user); err != nil {
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				logger.Debugf("user with telegram ID %d already exists", user.TelegramID)
				return "", "", ErrUserAlreadyExist
			}
		}
		logger.Errorf("failed to add user with telegram ID %d: %v", user.TelegramID, err)
		return "", "", ErrUserRegisrationFailed
	}
	logger.Infof("%s registered successfully", user)
	return a.generateTokens(user.ID)
}

func (a *AuthService) LoginWeb(ctx context.Context, name, password string) (string, string, error) {
	return "", "", nil
}

func (a *AuthService) LinkTelegramBot(ctx context.Context, bindToken string) error {
	return nil
}

func (a *AuthService) LoginTelegram(ctx context.Context, telegramID int64) (string, string, error) {
	user := &models.User{TelegramID: telegramID}
	if err := a.Repo.GetByTelegramID(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDBNoRows) {
			logger.Infof("user by telegram id %d not found", telegramID)
			return "", "", ErrUserNotFound
		}
		logger.Errorf("user with telegram id %v login failed: %v", telegramID, err)
		return "", "", fmt.Errorf("%w: %v", ErrLoginFailed, err)
	}

	return a.generateTokens(user.ID)
}

func (a *AuthService) LinkWeb(ctx context.Context, bindToken string) error {
	return nil
}

func (a *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := a.jwtManager.VerifyToken(refreshToken, jwt.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidRefreshToken) {
			logger.Infof("%v: %v", ErrRefreshTokenFailed, err)
			return "", "", ErrRefreshTokenFailed
		}
		return "", "", ErrInternalAuthService

	}

	return a.generateTokens(claims.ID)
}

// generateTokens return access, refresh tokens and optionally error
func (a *AuthService) generateTokens(userID string) (string, string, error) {
	accessToken, err := a.jwtManager.GenerateToken(userID, jwt.AccessToken)
	if err != nil {
		logger.Errorf("acсess token generation failed: %v", err)
		return "", "", ErrInternalAuthService

	}

	newRefreshToken, err := a.jwtManager.GenerateToken(userID, jwt.RefreshToken)
	if err != nil {
		logger.Errorf("refresh token generation failed: %v", err)
		return "", "", ErrInternalAuthService
	}
	return accessToken, newRefreshToken, nil
}
