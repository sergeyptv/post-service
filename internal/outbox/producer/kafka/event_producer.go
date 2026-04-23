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

	produce := func() error {
		err = k.producer.Kafka.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &k.producer.Config.Topic,
				Partition: kafka.PartitionAny,
			},
			Key:   []byte(event.UserUuid),
			Value: eventBytes,
		}, nil)

		if err != nil {
			return err
		}

		return nil
	}

	err = produce()
	if err != nil {
		retryTimeout := k.producer.Config.RetryTimeoutSec
		published := false

		for i := 1; i <= k.producer.Config.RetryAttempt; i++ {
			retryTimeout *= 2
			time.Sleep(time.Duration(retryTimeout) * time.Second)

			err = produce()
			if err == nil {
				published = true

				break
			}
		}

		if !published {
			topicDlq := k.producer.Config.Topic + ".dlq"

			err = k.producer.Kafka.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topicDlq,
					Partition: kafka.PartitionAny,
				},
				Key:   []byte(event.UserUuid),
				Value: eventBytes,
			}, nil)

			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	return nil
}
