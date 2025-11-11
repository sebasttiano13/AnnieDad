package domains

import (
	"encoding/json"
	"fmt"
	"time"
)

type FileStatus string

const (
	FileStatusNew       FileStatus = "new"
	FileStatusActive    FileStatus = "active"
	FileStatusUploading FileStatus = "uploading"
	FileStatusDeleted   FileStatus = "deleted"
)

type MimeType string

const (
	MimeJPEG MimeType = "image/jpeg"
	MimePNG  MimeType = "image/png"
	MimeGIF  MimeType = "image/gif"
	MimeMP4  MimeType = "video/mp4"
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
	// Group is a user groups
	Group struct {
		ID        string
		Name      string
		OwnerID   string
		IsShared  bool
		CreatedAt string
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
	// File is a file domain
	File struct {
		ID          string
		GroupID     string
		UploaderID  string
		FileName    string
		StoragePath string
		MimeType    MimeType
		SizeBytes   int64
		UploadedAt  time.Time
		Status      FileStatus
		IsShared    bool
		Meta        map[string]any
	}
)

func (r *RefreshToken) String() string {
	return fmt.Sprintf("<RefreshToken id=%s user_id=%s issued_at=%s expires_at=%s>", r.ID, r.UserID, r.IssuedAt, r.ExpiresAt)
}

func (f *File) MarshalMeta() ([]byte, error) {
	return json.Marshal(f.Meta)
}

func (f *File) UnmarshalMeta(data []byte) error {
	return json.Unmarshal(data, &f.Meta)
}
