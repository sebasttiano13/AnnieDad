package models

import (
	"database/sql"
	"fmt"
	"time"
)

type (
	// UserRecord is a client avatar in db
	UserRecord struct {
		ID             string
		TelegramID     int64          `db:"telegram_id"`
		Name           sql.NullString `db:"username"`
		HashedPassword sql.NullString `db:"password"`
		Email          sql.NullString `db:"email"`
		RegisteredAT   string         `db:"registered_at"`
	}
	// ApiClientRecord is an apps for clients in db
	ApiClientRecord struct {
		ID        string         `db:"id"`
		Name      sql.NullString `db:"name"`
		Token     sql.NullString `db:"token"`
		CreatedAt string         `db:"created_at"`
	}
	// RefreshTokenRecord is an api key stored in db
	RefreshTokenRecord struct {
		ID        string `db:"id"`
		Token     string `db:"token"`
		UserID    string `db:"user_id"`
		IssuedAt  string `db:"issued_at"`
		ExpiresAt string `db:"expires_at"`
		Revoked   bool   `db:"revoked"`
	}
	FileRecord struct {
		ID          string         `db:"id"`
		GroupID     sql.NullString `db:"group_id"`
		UploaderID  sql.NullString `db:"uploader_id"`
		FileName    string         `db:"file_name"`
		StoragePath string         `db:"storage_path"`
		MimeType    string         `db:"mime_type"`
		SizeBytes   int64          `db:"size_bytes"`
		UploadedAt  time.Time      `db:"uploaded_at"`
		Status      string         `db:"status"`
		IsShared    bool           `db:"is_shared"`
		Meta        []byte         `db:"meta"`
	}
)

func (u *UserRecord) String() string {
	return fmt.Sprintf("<UserRecord id=%s telegram=%d name=%s email=%s>", u.ID, u.TelegramID, u.Name.String, u.Email.String)
}
