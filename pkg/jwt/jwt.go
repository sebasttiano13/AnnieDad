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
	accessSecretKey      string
	refreshSecret        string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTManager(
	accessSecretKey string,
	refreshSecretKey string,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *JWTManager {
	return &JWTManager{accessSecretKey, refreshSecretKey, accessTokenDuration, refreshTokenDuration}
}

type Claims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

func (j *JWTManager) GenerateToken(id string, tokenType JWTTokenType) (string, error) {
	var errMessage = ErrJWTManager
	var expiration time.Time
	var secret string
	switch tokenType {
	case AccessToken:
		expiration = time.Now().Add(j.accessTokenDuration)
		secret = j.accessSecretKey
		errMessage = ErrAccessTokenCreationFailed
	case RefreshToken:
		expiration = time.Now().Add(j.refreshTokenDuration)
		secret = j.refreshSecret
		errMessage = ErrRefreshTokenCreationFailed
	default:
		return "", fmt.Errorf("%w: unknown token type", errMessage)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Subject:   id,
			Issuer:    "localhost:8080/api/user/login",
			Audience:  []string{"localhost:8080"},
		},
		ID: id,
	})
	newToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errMessage, err)
	}
	return newToken, nil
}

func (j *JWTManager) VerifyToken(checkToken string, tokenType JWTTokenType) (*Claims, error) {
	var errMessage = ErrJWTManager
	var secret string
	switch tokenType {
	case AccessToken:
		secret = j.accessSecretKey
		errMessage = ErrInvalidAccessToken
	case RefreshToken:
		secret = j.refreshSecret
		errMessage = ErrInvalidRefreshToken
	default:
		return nil, fmt.Errorf("%w: unknown token type", errMessage)
	}
	token, err := jwt.ParseWithClaims(
		checkToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrTokenMethod
			}

			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMessage, err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}
