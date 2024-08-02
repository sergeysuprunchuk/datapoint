package database

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Postgres = "postgres"
	Sqlite3  = "sqlite3"
)

var placeholders = map[string]sq.PlaceholderFormat{Postgres: sq.Dollar, Sqlite3: sq.Question}

const txKey = "tx"

type Database struct {
	*sql.DB
	B sq.StatementBuilderType
}

func New(driverName, dataSourceName string) (*Database, error) {
	db := new(Database)
	err := db.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) Open(driverName, dataSourceName string) error {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	if err = d.Ping(); err != nil {
		_ = d.Close()
		return err
	}

	if db.DB != nil {
		_ = db.DB.Close()
	}

	db.DB = d
	db.B = sq.StatementBuilder.
		RunWith(db).
		PlaceholderFormat(placeholders[driverName])

	return nil
}

func (db *Database) Close() {
	_ = db.DB.Close()
}

type handler func(context.Context) error

func (db *Database) ReadCommitted(ctx context.Context, h handler) error {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	if ok {
		return h(ctx)
	}

	var err error
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted}); err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %s", err)
	}

	ctx = context.WithValue(ctx, txKey, tx)

	if err = h(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *Database) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return db.DB.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return db.DB.QueryRowContext(ctx, query, args...)
}

func (db *Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return db.DB.ExecContext(ctx, query, args...)
}
