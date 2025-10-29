package models

import (
	"database/sql"
	"fmt"
)

type (
	// User is a client avatar
	User struct {
		ID             string
		TelegramID     int64          `db:"telegram_id"`
		Name           sql.NullString `json:"username" db:"username" valid:"required,type(string)"`
		HashedPassword string         `json:"password" db:"password" valid:"required,type(string)"`
		Email          sql.NullString `json:"email" db:"email" valid:"email"`
		RegisteredAT   string         `db:"registered_at"`
	}
	// ApiClient is an apps for clients
	ApiClient struct {
		ID        string
		Name      string
		Token     string
		CreatedAt string `db:"created_at"`
	}
)

func (u *User) String() string {
	return fmt.Sprintf("<User id=%s telegram=%d name=%s email=%s>", u.ID, u.TelegramID, u.Name, u.Email)
}
