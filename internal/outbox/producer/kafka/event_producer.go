package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sergeyptv/post_service/internal/outbox/domain"
	"github.com/sergeyptv/post_service/internal/platform/kafka_produce"
	"time"
)

type kafkaEventProducer struct {
	producer *kafka_produce.Producer
}

func NewKafkaEventProducer(p *kafka_produce.Producer) *kafkaEventProducer {
	return &kafkaEventProducer{
		producer: p,
	}
}

func (k *kafkaEventProducer) Publish(event domain.UserRegisteredEvent) error {
	const op = "producer.kafka.Publish"

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	deliveryChan := make(chan kafka.Event, 1)

	err = k.producer.Kafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.producer.Config.Topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(event.Uuid),
		Value: eventBytes,
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			return m.TopicPartition.Error
		}

	case <-time.After(5 * time.Second):
		return fmt.Errorf("%s: delivery timeout", op)
	}

	return nil
}
