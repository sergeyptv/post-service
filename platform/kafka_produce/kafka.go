package kafka_produce

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Config struct {
	Addr            string  `env:"ADDR" env-required`
	Topic           string  `env:"TOPIC" env-required`
	RetryAttempt    int     `env:"RETRY_ATTEMPTS" env-required`
	RetryTimeoutSec float64 `env:"RETRY_TIMEOUT_SEC" env-required`
}

type Producer struct {
	Config Config
	Kafka  *kafka.Producer
}

func NewProducer(ctx context.Context, c Config) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": c.Addr,
	})
	if err != nil {
		return nil, err
	}

	return &Producer{
		Kafka:  producer,
		Config: c,
	}, nil
}

func (p *Producer) Close() {
	p.Kafka.Flush(15 * 1000)
	p.Kafka.Close()
}
