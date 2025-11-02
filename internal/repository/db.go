package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano13/AnnieDad/internal/models"
)

var ErrDBTransaction = errors.New("sql: transaction error")
var ErrDB = errors.New("sql: db error")

// ErrDBNoRows error if no row was in db response.
var ErrDBNoRows = errors.New("sql: no rows in result set")

// DBStorage makes db queries
type DBStorage struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewDBStorage constructor for DBStorage
func NewDBStorage(d *sqlx.DB) *DBStorage {
	return &DBStorage{db: d, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

// GetByUsername returns user by username and hashed password
func (d *DBStorage) GetByUsername(ctx context.Context, user *models.User) error {
	query, args, err := d.psql.Select(UserID, UserTelegramID, UserName, UserPassword).From("users").Where(sq.Eq{UserName: user.Name}).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	if err := d.db.GetContext(ctx, user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDBNoRows
		}
		return fmt.Errorf("%w: %v", ErrDB, err)
	}
	return nil
}

// GetByTelegramID returns user by telegram id
func (d *DBStorage) GetByTelegramID(ctx context.Context, user *models.User) error {
	query, args, err := d.psql.Select(UserID, UserTelegramID, UserName).From("users").Where(sq.Eq{UserTelegramID: user.TelegramID}).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	if err := d.db.GetContext(ctx, user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDBNoRows
		}
		return fmt.Errorf("%w: %v", ErrDB, err)
	}
	return nil
}

// AddUser creates new user via username and password
func (d *DBStorage) AddUser(ctx context.Context, user *models.User) error {
	query, args, err := d.psql.Insert("users").
		Columns(UserName, UserPassword).
		Values(user.Name, user.HashedPassword).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBTransaction, err)
	}
	var id string
	if err := tx.GetContext(ctx, &id, query, args...); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", ErrDB, err)
	}
	user.ID = id

	tx.Commit()
	return nil
}

// AddBotUser creates new user via telegram id
func (d *DBStorage) AddBotUser(ctx context.Context, user *models.User) error {
	query, args, err := d.psql.Insert("users").
		Columns(UserTelegramID, UserName).
		Values(user.TelegramID, user.Name).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBTransaction, err)
	}
	var id string
	if err := tx.GetContext(ctx, &id, query, args...); err != nil {
		tx.Rollback()
		return err
	}
	user.ID = id

	tx.Commit()
	return nil
}

// LinkTelegramUser set telegram id to user
func (d *DBStorage) LinkTelegramUser(ctx context.Context, user *models.User) error {
	query, args, err := d.psql.Update("users").
		Set(UserTelegramID, user.TelegramID).
		Where(sq.Eq{UserID: user.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBTransaction, err)
	}
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", ErrDB, err)
	}

	tx.Commit()

	return nil
}
