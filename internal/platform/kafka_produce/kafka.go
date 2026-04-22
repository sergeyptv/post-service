package kafka_produce

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Config struct {
	Addr            []string `env:"ADDR" env-prefix:"KAFKA_PRODUCER_" env-required`
	Topic           string   `env:"TOPIC" env-prefix:"KAFKA_PRODUCER_" env-required`
	RetryAttempt    int      `env:"RETRY_ATTEMPTS" env-prefix:"KAFKA_PRODUCER_" env-required`
	RetryTimeoutSec float64  `env:"RETRY_TIMEOUT_SEC" env-prefix:"KAFKA_PRODUCER_" env-required`
}

type Producer struct {
	Kafka  *kafka.Producer
	Config Config
}

func NewProducer(ctx context.Context, c Config) (*Producer, chan string, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": c.Addr,
	})
	if err != nil {
		return nil, nil, err
	}

	eventDeliveryStatus := make(chan string)

	go func(ctx context.Context, eventDeliveryStatus chan string) {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					eventDeliveryStatus <- fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					eventDeliveryStatus <- fmt.Sprintf("Delivered message to %v\n", ev.TopicPartition)
				}
			}

			select {
			case <-ctx.Done():
				eventDeliveryStatus <- fmt.Sprintf("Kafka producer context closed: %v\n", ctx.Err())
				return

			default:
			}
		}
	}(ctx, eventDeliveryStatus)

	return &Producer{
		Kafka:  producer,
		Config: c,
	}, eventDeliveryStatus, nil
}

func (p *Producer) Close() {
	p.Kafka.Flush(15 * 1000)
	p.Kafka.Close()
}
