package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sergeyptv/post_service/platform/config"
	"github.com/sergeyptv/post_service/platform/kafka_produce"
	"github.com/sergeyptv/post_service/platform/postgres"
)

type Config struct {
	App                config.App           `env-prefix:"APP_"`
	WorkerFrequencySec int                  `env:"APP_WORKER_FREQUENCY_SEC" env-required`
	BatchSize          int                  `env:"APP_BATCH_SIZE" env-required`
	Postgres           postgres.Config      `env-prefix:"POSTGRES_"`
	KafkaProducer      kafka_produce.Config `env-prefix:"KAFKA_PRODUCER_"`
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
