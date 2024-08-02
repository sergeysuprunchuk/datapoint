package dbmodel

import (
	"context"
	"database/sql"
	"datapoint/pkg/database"
	"fmt"
	"github.com/google/uuid"
)

type DB struct {
	ID   string
	Info Info
	db   *database.Database
}

func New(info Info) (*DB, error) {
	driverName, dataSourceName := info.Config.Parse()
	db, err := database.New(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{ID: uuid.NewString(), Info: info, db: db}, err
}

func (db *DB) Open() error {
	driverName, dataSourceName := db.Info.Config.Parse()
	d, err := database.New(driverName, dataSourceName)
	db.db = d
	return err
}

func (db *DB) SetInfo(info Info) error {
	driverName, dataSourceName := info.Config.Parse()
	if err := db.db.Open(driverName, dataSourceName); err != nil {
		return err
	}
	db.Info = info
	return nil
}

func (db *DB) Close() {
	db.db.Close()
}

type Info struct {
	Name   string
	Config Config
}

const (
	PostgreSQL = "PostgreSQL"
)

var (
	drivers = map[string]string{PostgreSQL: database.Postgres}
	schemas = map[string]string{PostgreSQL: "postgresql"}
)

type Config struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Name     string
	Driver   string
}

func (c *Config) Parse() (string, string) {
	return drivers[c.Driver],
		fmt.Sprintf(
			"%s://%s:%s@%s:%d/%s?sslmode=disable",
			schemas[c.Driver], c.User, c.Password, c.Host, c.Port, c.Name,
		)
}

func (db *DB) Check() error {
	if db.db == nil {
		return db.Open()
	}
	return nil
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if err := db.Check(); err != nil {
		return nil, err
	}

	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if err := db.Check(); err != nil {
		return nil, err
	}

	return db.db.ExecContext(ctx, query, args...)
}
