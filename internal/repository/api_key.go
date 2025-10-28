package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/sebasttiano13/AnnieDad/internal/models"
)

type ApiKeyRepo interface {
	GetApiClient(ctx context.Context, client *models.ApiClient) error
}

type DBApiKeyChecker struct {
	repo ApiKeyRepo
}

func NewDBApiKeyChecker(repo ApiKeyRepo) *DBApiKeyChecker {
	return &DBApiKeyChecker{repo: repo}
}

func (c *DBApiKeyChecker) Validate(ctx context.Context, token string) (bool, error) {
	client := &models.ApiClient{Token: token}
	if err := c.repo.GetApiClient(ctx, client); err != nil {
		if errors.Is(err, ErrDBNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *DBStorage) GetApiClient(ctx context.Context, client *models.ApiClient) error {
	query, args, err := d.psql.Select(
		ApiClientID,
		ApiClientName,
		ApiClientToken,
		ApiClientCreatedAt).
		From("api_clients").Where(sq.Eq{ApiClientToken: client.Token}).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	if err := d.db.GetContext(ctx, client, query, args); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDBNoRows
		}
		return fmt.Errorf("%w: %v", ErrDB, err)
	}
	return nil
}
