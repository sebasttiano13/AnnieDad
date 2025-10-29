package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sebasttiano13/AnnieDad/internal/models"
	"github.com/sebasttiano13/AnnieDad/internal/repository"
	"github.com/sebasttiano13/AnnieDad/pkg/jwt"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
)

var ErrUserNotFound = errors.New("user not found")
var ErrLoginFailed = errors.New("login failed")
var ErrRefreshTokenFailed = errors.New("refresh login failed")

func (a *AuthService) Register(ctx context.Context, name, password string) (string, string, error) {
	return "", "", nil
}

func (a *AuthService) LoginWeb(ctx context.Context, name, password string) (string, string, error) {
	return "", "", nil
}

func (a *AuthService) LinkTelegramBot(ctx context.Context, bindToken string) error {
	return nil
}

func (a *AuthService) LoginBot(ctx context.Context, telegramID int64) (string, string, error) {
	user := &models.User{TelegramID: telegramID}
	if err := a.Repo.GetByTelegramID(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDBNoRows) {
			logger.Infof("user not found, creating new user for telegramID=%d", telegramID)
			user = &models.User{
				TelegramID: telegramID,
			}
			if err := a.Repo.AddBotUser(ctx, user); err != nil {
				logger.Errorf("failed to add bot user for telegramID=%d: %v", telegramID, err)
				return "", "", ErrInternalAuthService
			}
			logger.Infof("%s created successfully", user)
		} else {
			logger.Errorf("user with telegram id %v login failed: %v", telegramID, err)
			return "", "", fmt.Errorf("%w: %v", ErrLoginFailed, err)
		}
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
