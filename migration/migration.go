package migration

import (
	"database/sql"
	"os"
)

type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func FromFile(e Executor, name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	_, err = e.Exec(string(data))
	return err
}
