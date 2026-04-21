package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

type OutboxRepository interface {
	CreateEvent(ctx context.Context, event domain.UserRegisteredEvent) (string, error)
}
