package kafka_consume

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Config struct {
	Addr    string `env:"ADDR" env-required`
	Topic   string `env:"TOPIC" env-required`
	GroupId string `env:"GROUP_ID" env-required`
}

type Consumer struct {
	Config Config
	Kafka  *kafka.Consumer
	Msgs   chan *kafka.Message
}

func NewConsumer(ctx context.Context, c Config) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": c.Addr,
		"group.id":          c.GroupId,
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
