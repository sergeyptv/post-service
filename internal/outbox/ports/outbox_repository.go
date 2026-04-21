package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/outbox/domain"
)

type OutboxRepository interface {
	GetPending(ctx context.Context, limit int) ([]domain.UserRegisteredEvent, error)
	MarkSent(ctx context.Context, userUuid []string) error
	MarkFailed(ctx context.Context, userUuid []string) error
}
