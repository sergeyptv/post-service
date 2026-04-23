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
	db postgres.DBTX
}

func NewPostgresOutboxRepository(db postgres.DBTX) *postgresOutboxRepository {
	return &postgresOutboxRepository{
		db: db,
	}
}

func (p *postgresOutboxRepository) GetPending(ctx context.Context, tx postgres.DBTX, limit int) ([]domain.UserRegisteredEvent, error) {
	const op = "repository.postgres.GetPending"

	var eventUuids []string
	var event domain.UserRegisteredEvent
	var events []domain.UserRegisteredEvent

	rows, err := tx.Query(ctx,
		`UPDATE outbox.event SET status = 'processing', updated_at = now()
            	WHERE uuid IN (
            		SELECT uuid
            		FROM outbox.event
            		WHERE status = 'pending'
            			OR (status = 'processing' AND updated_at < now() - interval '30 seconds')
            		FOR UPDATE SKIP LOCKED
            		LIMIT $1
            	)
        		RETURNING uuid, version, user_uuid, user_email, registered_at`,
		limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrEventsNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&event.Uuid, &event.Version, &event.UserUuid, &event.UserEmail, &event.RegisteredAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		eventUuids = append(eventUuids, event.Uuid)
		events = append(events, event)
		event = domain.UserRegisteredEvent{}
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (p *postgresOutboxRepository) MarkSent(ctx context.Context, eventUuid string) error {
	const op = "repository.postgres.MarkSent"

	_, err := p.db.Exec(ctx,
		"UPDATE outbox.event SET status = 'sent' WHERE uuid = $1",
		eventUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *postgresOutboxRepository) MarkFailed(ctx context.Context, eventUuid string) error {
	const op = "repository.postgres.MarkFailed"

	_, err := p.db.Exec(ctx,
		"UPDATE outbox.event SET status = 'failed' WHERE uuid = $1",
		eventUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
