package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/outbox/domain"
	"github.com/sergeyptv/post_service/internal/outbox/repository"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
)

type postgresOutboxRepository struct {
	pool *postgres.Pool
}

func NewPostgresOutboxRepository(pool *postgres.Pool) *postgresOutboxRepository {
	return &postgresOutboxRepository{
		pool: pool,
	}
}

func (p *postgresOutboxRepository) GetPending(ctx context.Context, limit int) ([]domain.UserRegisteredEvent, error) {
	const op = "repository.postgres.GetPending"

	var event domain.UserRegisteredEvent
	var events []domain.UserRegisteredEvent

	rows, err := p.pool.Db.Query(ctx,
		"SELECT version, user_uuid, user_email, registered_at FROM outbox.event WHERE status = 'pending' LIMIT $1 FOR UPDATE SKIP LOCKED",
		limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrEventsNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&event.Version, &event.UserUuid, &event.UserEmail, &event.RegisteredAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		events = append(events, event)
		event = domain.UserRegisteredEvent{}
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (p *postgresOutboxRepository) MarkSent(ctx context.Context, userUuid []string) error {
	const op = "repository.postgres.MarkSent"

	_, err := p.pool.Db.Exec(ctx,
		"UPDATE outbox.event SET status = 'sent' WHERE user_uuid = $1",
		userUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *postgresOutboxRepository) MarkFailed(ctx context.Context, userUuid []string) error {
	const op = "repository.postgres.MarkFailed"

	_, err := p.pool.Db.Exec(ctx,
		"UPDATE outbox.event SET status = 'failed' WHERE user_uuid = $1",
		userUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
