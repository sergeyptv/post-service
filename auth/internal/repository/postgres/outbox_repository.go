package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/auth/internal/domain"
	"time"
)

type postgresOutboxRepository struct{}

func NewPostgresOutboxRepository() *postgresOutboxRepository {
	return &postgresOutboxRepository{}
}

func (p *postgresOutboxRepository) CreateEvent(ctx context.Context, tx pgx.Tx, event domain.UserRegisteredEvent) (string, error) {
	const op = "repository.postgres.CreateEvent"

	var eventUuid string

	err := tx.QueryRow(ctx,
		"INSERT INTO outbox.event (version, user_uuid, user_email, registered_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING uuid",
		event.Version, event.UserUuid, event.UserEmail, time.Now().UTC(), time.Now().UTC(),
	).Scan(&eventUuid)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return eventUuid, nil
}
