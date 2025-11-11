package repository

import "errors"

var ErrDBSqlBuilder = errors.New("sql: builder error")

// user table fields
const (
	UserID           = "id"
	UserTelegramID   = "telegram_id"
	UserName         = "username"
	UserPassword     = "password"
	UserEmail        = "email"
	UserRegisteredAt = "registered_at"
)

// group table fields
const (
	GroupID        = "id"
	GroupName      = "name"
	GroupOwnerID   = "owner_id"
	GroupIsShared  = "is_shared"
	GroupCreatedAt = "created_at"
)

// api_client table fields
const (
	ApiClientID        = "id"
	ApiClientName      = "name"
	ApiClientToken     = "token"
	ApiClientCreatedAt = "created_at"
)

// refresh_tokens table fields

const (
	RefreshID        = "id"
	RefreshUserID    = "user_id"
	RefreshToken     = "token"
	RefreshIssuedAt  = "issued_at"
	RefreshExpiresAt = "expires_at"
	RefreshRevoked   = "revoked"
)
