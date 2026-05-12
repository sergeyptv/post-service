package ports

import (
	"context"
	"github.com/sergeyptv/post_service/outbox/internal/domain"
	"github.com/sergeyptv/post_service/platform/postgres"
)

type OutboxRepository interface {
	GetPending(ctx context.Context, tx postgres.DBTX, limit int) ([]domain.UserRegisteredEvent, error)
	MarkSent(ctx context.Context, eventUuids []string) error
	MarkFailed(ctx context.Context, eventUuids []string) error
}
