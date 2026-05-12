package ports

import (
	"github.com/sergeyptv/post_service/outbox/internal/domain"
)

type EventPublisher interface {
	Publish(event domain.UserRegisteredEvent) error
}
