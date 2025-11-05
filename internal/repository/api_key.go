package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/sebasttiano13/AnnieDad/internal/domains"
	"github.com/sebasttiano13/AnnieDad/internal/models"
)

var ErrBadAPIKeyFormat = errors.New("bad api key format")

type ApiKeyRepo interface {
	GetApiClient(ctx context.Context, token string) (*domains.ApiClient, error)
}

type DBApiKeyChecker struct {
	repo ApiKeyRepo
}

func NewDBApiKeyChecker(repo ApiKeyRepo) *DBApiKeyChecker {
	return &DBApiKeyChecker{repo: repo}
}

func (c *DBApiKeyChecker) Validate(ctx context.Context, token string) (bool, error) {

	_, err := uuid.Parse(token)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrBadAPIKeyFormat, err)
	}

	_, err = c.repo.GetApiClient(ctx, token)
	if err != nil {
		if errors.Is(err, ErrDBNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *DBStorage) GetApiClient(ctx context.Context, token string) (*domains.ApiClient, error) {
	query, args, err := d.psql.Select(
		ApiClientID,
		ApiClientName,
		ApiClientToken,
		ApiClientCreatedAt).
		From("api_clients").Where(sq.Eq{ApiClientToken: token}).ToSql()

	client := &models.ApiClientRecord{}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	if err := d.db.GetContext(ctx, client, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDBNoRows
		}
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}
	return apiClientRecordToDomain(client), nil
}
