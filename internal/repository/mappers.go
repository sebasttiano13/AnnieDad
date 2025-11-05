package repository

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/sebasttiano13/AnnieDad/internal/domains"
	"github.com/sebasttiano13/AnnieDad/internal/models"
)

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
	// turn timestamp to time.Time
	issuedAtsec, err := strconv.ParseInt(r.IssuedAt, 10, 64)
	if err != nil {
		return nil, err
	}
	issuedAtts := time.Unix(issuedAtsec, 0)
	expiresAtsec, err := strconv.ParseInt(r.IssuedAt, 10, 64)
	if err != nil {
		return nil, err
	}
	expiresAtts := time.Unix(expiresAtsec, 0)

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
