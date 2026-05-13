package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sergeyptv/post_service/platform/config"
	"github.com/sergeyptv/post_service/platform/kafka_consume"
	"github.com/sergeyptv/post_service/platform/migrator"
	"github.com/sergeyptv/post_service/platform/postgres"
)

type Config struct {
	App           config.App
	Migrator      migrator.Config
	Postgres      postgres.Config
	KafkaConsumer kafka_consume.Config
}

func MustLoad() *Config {
	return mustParseEnv()
}

func mustParseEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
