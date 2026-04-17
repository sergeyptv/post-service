package domain

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, topic string, event UserRegisteredEvent) error
}
