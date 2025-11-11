package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/sebasttiano13/AnnieDad/internal/domains"
	"github.com/sebasttiano13/AnnieDad/internal/models"
)

func RecordToDomainFile(r *models.FileRecord) (*domains.File, error) {
	meta := make(map[string]any)
	if len(r.Meta) > 0 {
		if err := json.Unmarshal(r.Meta, &meta); err != nil {
			return nil, err
		}
	}

	return &domains.File{
		ID:          r.ID,
		GroupID:     r.GroupID.String,
		UploaderID:  r.UploaderID.String,
		FileName:    r.FileName,
		StoragePath: r.StoragePath,
		MimeType:    domains.MimeType(r.MimeType),
		SizeBytes:   r.SizeBytes,
		UploadedAt:  r.UploadedAt,
		Status:      domains.FileStatus(r.Status),
		IsShared:    r.IsShared,
		Meta:        meta,
	}, nil
}

func userRecordToDomain(r *models.UserRecord) *domains.User {
	if r == nil {
		return nil
	}
	return &domains.User{
		ID:           r.ID,
		TelegramID:   r.TelegramID,
		Name:         nullStringToString(r.Name),
		Email:        nullStringToString(r.Email),
		RegisteredAT: r.RegisteredAT,
	}
}

func userDomainToRecord(u *domains.User) *models.UserRecord {
	return &models.UserRecord{
		ID:           u.ID,
		TelegramID:   u.TelegramID,
		Name:         stringToNullString(u.Name),
		Email:        stringToNullString(u.Email),
		RegisteredAT: u.RegisteredAT,
	}
}

func refreshRecordToDomain(r *models.RefreshTokenRecord) (*domains.RefreshToken, error) {
	if r == nil {
		return nil, nil
	}

	issuedAtts, err := time.Parse(time.RFC3339, r.IssuedAt)
	if err != nil {
		return nil, err
	}

	expiresAtts, err := time.Parse(time.RFC3339, r.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &domains.RefreshToken{
		ID:        r.ID,
		UserID:    r.UserID,
		IssuedAt:  issuedAtts,
		ExpiresAt: expiresAtts,
		Revoked:   r.Revoked,
		Token:     r.Token,
	}, err
}

func apiClientRecordToDomain(r *models.ApiClientRecord) *domains.ApiClient {
	if r == nil {
		return nil
	}
	return &domains.ApiClient{
		ID:        r.ID,
		Name:      nullStringToString(r.Name),
		Token:     nullStringToString(r.Token),
		CreatedAt: r.CreatedAt,
	}
}

func apiClientDomainToRecord(r *domains.ApiClient) *models.ApiClientRecord {
	return &models.ApiClientRecord{
		ID:        r.ID,
		Name:      stringToNullString(r.Name),
		Token:     stringToNullString(r.Token),
		CreatedAt: r.CreatedAt,
	}
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
