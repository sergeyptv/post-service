package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/sergeyptv/post_service/internal/notification/domain"
	"github.com/sergeyptv/post_service/internal/platform/kafka_consume"
	"time"
)

type kafkaEventConsumer struct {
	consumer *kafka_consume.Consumer
}

func NewKafkaEventConsumer(consumer *kafka_consume.Consumer) *kafkaEventConsumer {
	return &kafkaEventConsumer{
		consumer: consumer,
	}
}

func (k *kafkaEventConsumer) Poll() (domain.UserRegisteredEvent, error) {
	const op = "consumer.kafka.Pull"

	var userRegisteredEvent domain.UserRegisteredEvent

	msg, err := k.consumer.Kafka.ReadMessage(time.Second)
	if err != nil {
		return domain.UserRegisteredEvent{}, fmt.Errorf("%s: %w", op, err)
	}

	err = json.Unmarshal(msg.Value, &userRegisteredEvent)
	if err != nil {
		return domain.UserRegisteredEvent{}, fmt.Errorf("%s: %w", op, err)
	}

	return userRegisteredEvent, nil
}

func (k *kafkaEventConsumer) Commit() error {
	const op = "consumer.kafka.Commit"

	_, err := k.consumer.Kafka.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
