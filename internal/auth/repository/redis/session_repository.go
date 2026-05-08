package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
)

type postgresSessionRepository struct {
	pool *postgres.Pool
}

func NewPostgresSessionRepository(pool *postgres.Pool) *postgresSessionRepository {
	return &postgresSessionRepository{
		pool: pool,
	}
}

func (p *postgresSessionRepository) CreateToken(ctx context.Context, userUuid string, token string) (string, error) {
	const op = "repository.postgres.CreateToken"

	var jti string

	err := p.pool.Db.QueryRow(ctx,
		`INSERT INTO token.storage (user_uuid, token)
				VALUES ($1, $2)
				RETURNING jti`,
		userUuid, token).Scan(&jti)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return jti, nil
}
func (p *postgresSessionRepository) GetToken(ctx context.Context, userUuid string) (string, error) {
	const op = "repository.postgres.GetToken"

	var jti string

	err := p.pool.Db.QueryRow(ctx,
		`SELECT jti
				FROM token.storage
				WHERE user_uuid = $1`,
		userUuid).Scan(&jti)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, repository.ErrTokenNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return jti, nil
}
func (p *postgresSessionRepository) UpdateToken(ctx context.Context, jti string, newToken string) error {
	const op = "repository.postgres.UpdateToken"

	cmdTag, err := p.pool.Db.Exec(ctx,
		`UPDATE token.storage
				SET token = $1
				WHERE jti = $2`,
		newToken, jti)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() < 1 {
		return fmt.Errorf("%s: %w", op, repository.ErrNoRowsAffected)
	}

	return nil
}
