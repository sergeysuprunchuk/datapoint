package dbrepo

import (
	"context"
	"datapoint/internal/model/dbmodel"
	"datapoint/internal/service/dbservice"
	"datapoint/pkg/database"
)

type repo struct {
	db *database.Database
}

var _ dbservice.DBRepo = (*repo)(nil)

var columns = []string{
	"id",
	"name",
	"host",
	"port",
	"db_user",
	"password",
	"db_name",
	"driver",
}

func (r *repo) GetList(ctx context.Context) ([]*dbmodel.DB, error) {
	rows, err := r.db.B.
		Select(columns...).
		From("database").
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []*dbmodel.DB
	for rows.Next() {
		d := new(dbmodel.DB)

		if err = rows.Scan(
			d.ID,
			d.Info.Name,
			d.Info.Config.Host,
			d.Info.Config.Port,
			d.Info.Config.User,
			d.Info.Config.Password,
			d.Info.Config.Name,
			d.Info.Config.Driver,
		); err != nil {
			return nil, err
		}

		list = append(list, d)
	}

	return list, nil
}

func (r *repo) Add(ctx context.Context, d dbmodel.DB) error {
	_, err := r.db.B.
		Insert("database").
		Columns(columns...).
		Values(
			d.ID,
			d.Info.Name,
			d.Info.Config.Host,
			d.Info.Config.Port,
			d.Info.Config.User,
			d.Info.Config.Password,
			d.Info.Config.Name,
			d.Info.Config.Driver,
		).
		ExecContext(ctx)
	return err
}

func (r *repo) Edit(ctx context.Context, d dbmodel.DB) error {
	_, err := r.db.B.
		Update("database").
		Set("name", d.Info.Name).
		Set("host", d.Info.Config.Host).
		Set("port", d.Info.Config.Port).
		Set("db_user", d.Info.Config.User).
		Set("password", d.Info.Config.Password).
		Set("db_name", d.Info.Config.Name).
		Set("driver", d.Info.Config.Driver).
		Where("id = ?", d.ID).
		ExecContext(ctx)
	return err
}

func (r *repo) Delete(ctx context.Context, id string) error {
	_, err := r.db.B.
		Delete("database").
		Where("id = ?", id).
		ExecContext(ctx)
	return err
}

func New(db *database.Database) *repo {
	return &repo{db: db}
}
