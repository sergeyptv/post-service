package ports

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

type OutboxRepository interface {
	CreateEvent(ctx context.Context, tx pgx.Tx, event domain.UserRegisteredEvent) (string, error)
}
