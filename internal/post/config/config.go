package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sergeyptv/post_service/internal/platform/config"
	"github.com/sergeyptv/post_service/internal/platform/httpserver"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
)

type Config struct {
	App      config.App
	Server   httpserver.Config
	Postgres postgres.Config
}

func MustLoad() *Config {
	return mustParseEnv()
}

func mustParseEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot get all env")
	}

	return &cfg
}
