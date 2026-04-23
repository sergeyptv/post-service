package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/notification/repository"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
)

type postgresNotificationRepository struct {
	db postgres.DBTX
}

func NewPostgresNotificationRepository(db postgres.DBTX) *postgresNotificationRepository {
	return &postgresNotificationRepository{
		db: db,
	}
}

func (p *postgresNotificationRepository) TryProcess(ctx context.Context, eventUuid string) error {
	const op = "repository.postgres.SetProcessing"

	var status string

	err := p.db.QueryRow(ctx,
		`INSERT INTO notification.processed_event (uuid, status, updated_at) 
				VALUES ($1, $2, now())
				ON CONFLICT (uuid) DO UPDATE SET
					status = EXCLUDED.status,
					updated_at = now()
				WHERE notification.processed_event.status = 'processing'
				  AND notification.processed_event.updated_at < now() - interval '60 seconds'
				RETURNING status`,
		eventUuid, "processing").Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			qerr := p.db.QueryRow(ctx,
				"SELECT status FROM notification.processed_event WHERE uuid = $1",
				eventUuid).Scan(&status)
			if qerr != nil {
				return fmt.Errorf("%s: %w", op, qerr)
			}

			if status == "success" {
				return fmt.Errorf("%s: %w", op, repository.ErrEventAlreadySuccess)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *postgresNotificationRepository) MarkSuccess(ctx context.Context, eventUuid string) error {
	const op = "repository.postgres.SetProcessed"

	_, err := p.db.Exec(ctx,
		"UPDATE notification.processed_event SET status = $1, updated_at = now() WHERE uuid = $2",
		"success", eventUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
