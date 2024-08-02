package dbmodel

import (
	"context"
	"database/sql"
	"datapoint/pkg/database"
	"fmt"
	sq "github.com/Masterminds/squirrel"
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

type FK struct {
	TableName  string
	ColumnName string
}

type Column struct {
	Name       string
	Type       string
	IsRequired bool
	IsPK       bool
	FK         *FK
}

type Table struct {
	Name       string
	ColumnList []*Column
}

func (db *DB) tableList(ctx context.Context, where sq.Sqlizer) ([]*Table, error) {
	rows, err := db.db.B.Select(
		"c.table_name",
		"c.column_name",
		"c.data_type",
		"c.is_nullable = 'NO' AND c.column_default IS NULL",
		"tc.constraint_type",
		"kcu2.table_name",
		"kcu2.column_name",
	).From("information_schema.columns c").
		LeftJoin("information_schema.key_column_usage kcu USING (table_name, column_name)").
		LeftJoin("information_schema.table_constraints tc USING (constraint_name)").
		LeftJoin("information_schema.referential_constraints rc USING (constraint_name)").
		LeftJoin("information_schema.key_column_usage kcu2 ON rc.unique_constraint_name = kcu2.constraint_name").
		Where("c.table_schema = 'public'").
		Where(where).
		OrderBy("c.table_name", "c.ordinal_position").
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var (
		tableList []*Table
		lastT     *Table
		lastC     *Column
	)

	for rows.Next() {
		var (
			t                    = new(Table)
			c                    = new(Column)
			constraint, fkT, fkC *string
		)

		if err = rows.Scan(&t.Name, &c.Name, &c.Type, &c.IsRequired, &constraint, &fkT, &fkC); err != nil {
			return nil, err
		}

		if lastT == nil || lastT.Name != t.Name {
			lastT = t
			tableList = append(tableList, t)
		}

		if lastC == nil || lastC.Name != c.Name {
			lastC = c
			lastT.ColumnList = append(lastT.ColumnList, c)
		}

		if constraint != nil && *constraint == "PRIMARY KEY" {
			lastC.IsPK = true
		} else if constraint != nil && *constraint == "FOREIGN KEY" && fkT != nil && fkC != nil {
			lastC.FK = &FK{TableName: *fkT, ColumnName: *fkC}
		}
	}

	return tableList, nil
}

func (db *DB) TableList(ctx context.Context) ([]*Table, error) {
	return db.tableList(ctx, nil)
}

func (db *DB) TableByName(ctx context.Context, name string) (*Table, error) {
	tableList, err := db.tableList(ctx, sq.Eq{"c.table_name": name})
	if err != nil {
		return nil, err
	}

	if len(tableList) == 0 {
		return nil, fmt.Errorf("таблицы с именем %s не существует", name)
	}

	return tableList[0], nil
}
