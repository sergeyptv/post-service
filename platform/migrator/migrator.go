package migrator

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Config struct {
	MigrationsDir string `env:"DIR" env-prefix:"MIGRATIONS_" env-required`
}

func Up(dir, dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute)

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err == nil {
			break
		}

		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return err
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	return goose.Up(db, dir)
}
