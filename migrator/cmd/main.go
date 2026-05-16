package main

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sergeyptv/post_service/platform/config"
	"github.com/sergeyptv/post_service/platform/logger"
	"github.com/sergeyptv/post_service/platform/migrator"
	"github.com/sergeyptv/post_service/platform/postgres"
)

type Config struct {
	App      config.App      `env-prefix:"APP_"`
	Migrator migrator.Config `env-prefix:"MIGRATIONS_"`
	Postgres postgres.Config `env-prefix:"POSTGRES_"`
}

func mustLoad() *Config {
	return mustParseEnv()
}

func mustParseEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func main() {
	cfg := mustLoad()

	log := logger.SetupLogger(cfg.App.Env)

	if err := runMigrations(cfg); err != nil {
		log.Error("Failed to up migrations", logger.Error(err))
		os.Exit(1)
	}
}

func runMigrations(cfg *Config) error {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	err := migrator.Up(cfg.Migrator.MigrationsDir, dsn)
	if err != nil {
		return err
	}

	return nil
}
