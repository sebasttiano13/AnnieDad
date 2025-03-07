package repository

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano13/AnnieDad/internal/models"
)

// ErrDBNoRows error if no row was in db response.
var ErrDBNoRows = errors.New("sql: no rows in result set")

// DBStorage makes db queries
type DBStorage struct {
	db *sqlx.DB
}

// NewDBStorage constructor for DBStorage
func NewDBStorage(d *sqlx.DB) *DBStorage {
	return &DBStorage{db: d}
}

func (d *DBStorage) GetUser(ctx context.Context, user *models.User) error {
	return nil
}

func (d *DBStorage) AddUser(ctx context.Context, user *models.User) error {
	return nil
}
