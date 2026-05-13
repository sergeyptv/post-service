package migrator

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Config struct {
	Dir string `env:"DIR" env-prefix:"MIGRATIONS_" env-required`
}

func Up(dir, dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	return goose.Up(db, dir)
}
