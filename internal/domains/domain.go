package domains

import (
	"fmt"
	"time"
)

type (
	// User is a client avatar domain
	User struct {
		ID             string
		TelegramID     int64
		Name           string
		HashedPassword string
		Email          string
		RegisteredAT   string
	}
	// ApiClient is an apps for clients domain
	ApiClient struct {
		ID        string
		Name      string
		Token     string
		CreatedAt string
	}
	// RefreshToken is a refresh token domain
	RefreshToken struct {
		ID        string
		Token     string
		UserID    string
		IssuedAt  time.Time
		ExpiresAt time.Time
		Revoked   bool
	}
)

func (r *RefreshToken) String() string {
	return fmt.Sprintf("<RefreshToken id=%s user_id=%s issued_at=%s expires_at=%s>", r.ID, r.UserID, r.IssuedAt, r.ExpiresAt)
}
