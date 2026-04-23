package ports

import (
	"github.com/sergeyptv/post_service/internal/outbox/domain"
)

type EventPublisher interface {
	Publish(event domain.UserRegisteredEvent) error
}
