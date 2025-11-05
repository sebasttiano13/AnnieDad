package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenType string

const (
	AccessToken  JWTTokenType = "access_token"
	RefreshToken JWTTokenType = "refresh_token"
)

var (
	ErrJWTManager                 = errors.New("jwt manager error")
	ErrAccessTokenCreationFailed  = errors.New("access token creation failed")
	ErrRefreshTokenCreationFailed = errors.New("refresh token creation failed")
	ErrInvalidAccessToken         = errors.New("invalid access token")
	ErrInvalidRefreshToken        = errors.New("invalid refresh token")
	ErrInvalidTokenClaims         = errors.New("invalid token claims")
	ErrTokenMethod                = errors.New("unexpected token signing method")
)

type JWTManager struct {
	accessSecretKey  string
	refreshSecretKey string
}

func NewJWTManager(
	accessSecretKey string,
	refreshSecretKey string,
) *JWTManager {
	return &JWTManager{accessSecretKey, refreshSecretKey}
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	ID      string `json:"id"`
	TokenID string `json:"token_id"`
}

func (j *JWTManager) GenerateAccessToken(userID string, duration time.Duration) (string, error) {
	expiration := time.Now().Add(duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Subject:   userID,
			Issuer:    "localhost:8080/api/user/login",
			Audience:  []string{"localhost:8080"},
		},
		ID: userID,
	})
	newToken, err := token.SignedString([]byte(j.accessSecretKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAccessTokenCreationFailed, err)
	}
	return newToken, nil
}

func (j *JWTManager) GenerateRefreshToken(userID, tokenID string, duration time.Duration) (string, error) {
	expiration := time.Now().Add(duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Subject:   userID,
			Issuer:    "localhost:8080/api/user/login",
			Audience:  []string{"localhost:8080"},
		},
		ID:      userID,
		TokenID: tokenID,
	})
	newToken, err := token.SignedString([]byte(j.refreshSecretKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
	}
	return newToken, nil
}

func (j *JWTManager) VerifyAccessToken(checkToken string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		checkToken,
		&AccessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrTokenMethod
			}

			return []byte(j.accessSecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidAccessToken, err)
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}

func (j *JWTManager) VerifyRefreshToken(checkToken string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		checkToken,
		&RefreshTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrTokenMethod
			}

			return []byte(j.refreshSecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}
