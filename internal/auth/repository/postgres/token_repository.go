package postgres

import (
	"context"
	"fmt"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
)

type postgresTokenRepository struct {
	pool *postgres.Pool
}

func NewPostgresTokenRepository(pool *postgres.Pool) *postgresTokenRepository {
	return &postgresTokenRepository{
		pool: pool,
	}
}

func (p *postgresTokenRepository) CreateToken(ctx context.Context, userUuid string, token string) (string, error) {
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
func (p *postgresTokenRepository) GetToken(ctx context.Context, userUuid string) (string, error) {
	const op = "repository.postgres.GetToken"

	var token string

	err := p.pool.Db.QueryRow(ctx,
		`SELECT token
				FROM token.storage
				WHERE user_uuid = $1`,
		userUuid).Scan(&token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
func (p *postgresTokenRepository) UpdateToken(ctx context.Context, jti string, newToken string) error {
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
		return repository.ErrNoRowsAffected
	}

	return nil
}
