package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano13/AnnieDad/internal/domains"
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
func (d *DBStorage) GetByUsername(ctx context.Context, userName string) (*domains.User, error) {
	query, args, err := d.psql.Select(UserID, UserTelegramID, UserName, UserPassword).From("users").Where(sq.Eq{UserName: userName}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	user := &models.UserRecord{}
	if err := d.db.GetContext(ctx, user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDBNoRows
		}
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}
	return userRecordToDomain(user), nil
}

// GetByTelegramID returns user by telegram id
func (d *DBStorage) GetByTelegramID(ctx context.Context, telegramID int64) (*domains.User, error) {
	query, args, err := d.psql.Select(UserID, UserTelegramID, UserName).From("users").Where(sq.Eq{UserTelegramID: telegramID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	user := &models.UserRecord{}
	if err := d.db.GetContext(ctx, user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDBNoRows
		}
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}
	return userRecordToDomain(user), nil
}

// AddUser creates new user via username and password
func (d *DBStorage) AddUser(ctx context.Context, userName string, userPassword string) (*domains.User, error) {
	query, args, err := d.psql.Insert("users").
		Columns(UserName, UserPassword).
		Values(UserPassword, userPassword).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBTransaction, err)
	}
	var id string
	if err := tx.GetContext(ctx, &id, query, args...); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}

	tx.Commit()
	user := &domains.User{ID: id}
	return user, nil
}

// AddBotUser creates new user via telegram id
func (d *DBStorage) AddBotUser(ctx context.Context, telegramID int64, userName string) (*domains.User, error) {
	query, args, err := d.psql.Insert("users").
		Columns(UserTelegramID, UserName).
		Values(telegramID, userName).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBTransaction, err)
	}
	var id string
	if err := tx.GetContext(ctx, &id, query, args...); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	user := &domains.User{ID: id}
	return user, nil
}

// LinkTelegramUser set telegram id to user
func (d *DBStorage) LinkTelegramUser(ctx context.Context, user *domains.User) error {
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

func (d *DBStorage) SaveRefresh(ctx context.Context, userID, token, tokenID string, expiresAt time.Time) error {
	query, args, err := d.psql.Insert("refresh_tokens").
		Columns(RefreshID, RefreshUserID, RefreshToken, RefreshExpiresAt).
		Values(tokenID, userID, token, expiresAt).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBTransaction, err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", ErrDB, err)
	}

	tx.Commit()
	return nil
}

func (d *DBStorage) GetRefresh(ctx context.Context, id string) (*domains.RefreshToken, error) {
	query, args, err := d.psql.Select(RefreshID, RefreshToken, RefreshUserID, RefreshIssuedAt, RefreshExpiresAt, RefreshRevoked).From("refresh_tokens").
		Where(sq.Eq{RefreshID: id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}
	refresh := &models.RefreshTokenRecord{}
	if err := d.db.GetContext(ctx, refresh, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDBNoRows
		}
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}
	return refreshRecordToDomain(refresh)
}

func (d *DBStorage) DeleteRefresh(ctx context.Context, id string) error {
	return nil
}

func (d *DBStorage) DeleteExpiredRefresh(ctx context.Context) (int64, error) {
	query, args, err := d.psql.Delete("refresh_tokens").Where(sq.Lt{RefreshExpiresAt: time.Now().UTC()}).ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDBSqlBuilder, err)
	}

	res, err := d.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDB, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDB, err)
	}

	return count, nil
}
