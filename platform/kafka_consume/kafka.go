package kafka_consume

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Config struct {
	Addr  []string `env:"ADDR" env-prefix:"KAFKA_CONSUMER_" env-required`
	Topic string   `env:"TOPIC" env-prefix:"KAFKA_CONSUMER_" env-required`
}

type Consumer struct {
	Config Config
	Kafka  *kafka.Consumer
	Msgs   chan *kafka.Message
}

func NewConsumer(ctx context.Context, c Config) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": c.Addr,
	})
	if err != nil {
		return nil, err
	}

	err = consumer.Subscribe(c.Topic, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Kafka:  consumer,
		Config: c,
		Msgs:   make(chan *kafka.Message, 0),
	}, nil
}

func (p *Consumer) Close() {
	if err := p.Kafka.Close(); err != nil {
		fmt.Printf("Error close kafka_consume consumer: %s\n", err)
	}
}
