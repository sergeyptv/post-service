package kafka_consume

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type Config struct {
	Addr  []string `env:"ADDR" env-prefix:"KAFKA_CONSUMER_" env-required`
	Topic string   `env:"TOPIC" env-prefix:"KAFKA_CONSUMER_" env-required`
}

type Consumer struct {
	Kafka  *kafka.Consumer
	Config Config
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

func (p *Consumer) Consume(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			msg, err := p.Kafka.ReadMessage(time.Second)
			if err == nil {
				p.Msgs <- msg
			}
		}
	}
}
