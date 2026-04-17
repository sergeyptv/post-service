package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/platform/kafka_produce"
)

type kafkaEventProducer struct {
	producer *kafka_produce.Producer
}

func NewKafkaEventProducer(p *kafka_produce.Producer) *kafkaEventProducer {
	return &kafkaEventProducer{
		producer: p,
	}
}

func (k *kafkaEventProducer) Publish(ctx context.Context, topic string, event domain.UserRegisteredEvent) error {
	const op = "repository.kafka.Publish"

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = k.producer.Kafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(event.UserUuid),
		Value: eventBytes,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
