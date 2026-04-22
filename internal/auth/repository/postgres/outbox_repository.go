package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/auth/repository"
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

func (p *postgresOutboxRepository) CreateEvent(ctx context.Context, tx pgx.Tx, event domain.UserRegisteredEvent) (string, error) {
	const op = "repository.postgres.CreateEvent"

	var eventUuid string

	err := tx.QueryRow(ctx,
		"INSERT INTO outbox.event (version, user_uuid, username, user_email, registered_at) VALUES ($1, $2, $3, $4, $5) RETURNING uuid",
		event.Version, event.UserUuid, event.Username, event.UserEmail, event.RegisteredAt,
	).Scan(&eventUuid)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return eventUuid, nil
}
