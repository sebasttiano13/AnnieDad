package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sebasttiano13/AnnieDad/internal/repository"
	"github.com/sebasttiano13/AnnieDad/pkg/jwt"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
)

var (
	ErrRefreshTokenFailed = errors.New("refresh login failed")
	ErrTokenService       = errors.New("token service internal error")
)

func (t *TokenService) GetTokens(ctx context.Context, userID string) (string, string, error) {
	accessToken, err := t.generateAccessToken(userID)
	if err != nil {
		return "", "", err

	}
	newRefreshToken, err := t.getRefreshToken(ctx, "", userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (t *TokenService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := t.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidRefreshToken) {
			logger.Infof("%v: %v", ErrRefreshTokenFailed, err)
			return "", "", ErrRefreshTokenFailed
		}
		return "", "", ErrTokenService
	}

	accessToken, err := t.generateAccessToken(claims.ID)
	if err != nil {
		return "", "", err

	}

	updatedRefreshToken, err := t.getRefreshToken(ctx, claims.TokenID, claims.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, updatedRefreshToken, err
}

func (t *TokenService) RevokeToken(ctx context.Context, refreshToken string) error {
	return nil
}

func (t *TokenService) generateAccessToken(userID string) (string, error) {
	accessToken, err := t.jwtManager.GenerateAccessToken(userID, t.opts.AccessTokenDuration)
	if err != nil {
		logger.Errorf("acсess token generation failed: %v", err)
		return "", ErrTokenService
	}
	return accessToken, nil
}

func (t *TokenService) generateRefreshToken(ctx context.Context, userID string) (string, error) {

	// Save to db
	expiration := time.Now().UTC().Add(t.opts.RefreshTokenDuration)

	newTokenid := uuid.New().String()

	refreshToken, err := t.jwtManager.GenerateRefreshToken(userID, newTokenid, t.opts.RefreshTokenDuration)
	if err != nil {
		logger.Errorf("refresh token generation failed: %v", err)
		return "", ErrTokenService
	}
	if err := t.repo.SaveRefresh(ctx, userID, refreshToken, newTokenid, expiration); err != nil {
		logger.Errorf("save refresh token failed: %v", err)
		return "", ErrTokenService
	}

	return refreshToken, nil
}

func (t *TokenService) getRefreshToken(ctx context.Context, refreshID, userID string) (string, error) {

	if refreshID == "" {
		return t.generateRefreshToken(ctx, userID)
	}
	// Check in db
	currentRefresh, err := t.repo.GetRefresh(ctx, refreshID)
	if err != nil {
		if errors.Is(err, repository.ErrDBNoRows) {
			logger.Infof("refresh by id %s not found. generating new", refreshID)
			return t.generateRefreshToken(ctx, userID)
		}
		return "", fmt.Errorf("%w: %v", ErrTokenService, err)
	}
	if time.Until(currentRefresh.ExpiresAt) < t.opts.RefreshTokenRenewBefore {
		logger.Infof("refresh by id %s soon will be expired. Renew refresh token", refreshID)
		return t.generateRefreshToken(ctx, userID)
	}

	return currentRefresh.Token, nil
}

func (t *TokenService) RunCleanup(ctx context.Context) {
	ticker := time.NewTicker(t.opts.RefreshCleanupInterval)
	defer ticker.Stop()

	logger.Infof("refresh cleanup started with frequency %s", t.opts.RefreshCleanupInterval)

	for {
		select {
		case <-ticker.C:
			t.cleanup(ctx)
		case <-ctx.Done():
			logger.Info("refresh cleanup stopped")
			return
		}
	}
}

func (t *TokenService) cleanup(ctx context.Context) {
	deleted, err := t.repo.DeleteExpiredRefresh(ctx)
	if err != nil {
		logger.Errorf("failed to delete expired refresh tokens: %v", err)
		return
	}
	if deleted > 0 {
		logger.Debugf("deleted %d expired refresh tokens", deleted)
	}
}
