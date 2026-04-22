package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/outbox/domain"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
)

type OutboxRepository interface {
	GetPending(ctx context.Context, tx postgres.DBTX, limit int) ([]string, []domain.UserRegisteredEvent, error)
	MarkSent(ctx context.Context, eventUuids []string) error
	MarkFailed(ctx context.Context, eventUuids []string) error
}
