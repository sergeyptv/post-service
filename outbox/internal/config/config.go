package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sergeyptv/post_service/platform/config"
	"github.com/sergeyptv/post_service/platform/kafka_produce"
	"github.com/sergeyptv/post_service/platform/postgres"
)

type Config struct {
	App                config.App
	WorkerFrequencySec int `env:"WORKER_FREQUENCY_SEC" env-prefix:"APP_" env-required`
	BatchSize          int `env:"BATCH_SIZE" env-prefix:"APP_" env-required`
	Postgres           postgres.Config
	KafkaProducer      kafka_produce.Config
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
