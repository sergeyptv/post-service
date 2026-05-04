package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/auth/domain"
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

	var tokenUuid string

	err := p.pool.Db.QueryRow(ctx,
		`UPDATE token.storage
				SET token = $1
				WHERE user_uuid = $2
				RETURNING uuid`,
		token, userUuid).Scan(&tokenUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			qerr := p.pool.Db.QueryRow(ctx,
				`INSERT INTO token.storage (user_uuid, token)
						VALUES ($1, $2)
						RETURNING uuid`,
				userUuid, token).Scan(&tokenUuid)

			if qerr != nil {
				return "", fmt.Errorf("%s: %w", op, qerr)
			}
		} else {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	return tokenUuid, nil
}
func (p *postgresTokenRepository) GetToken(ctx context.Context, tokenUuid string) (domain.Token, error) {
}
func (p *postgresTokenRepository) UpdateToken(ctx context.Context, tokenUuid string, updToken domain.UpdateToken) error {
}
func (p *postgresTokenRepository) DeleteToken(ctx context.Context, tokenUuid string) error {}
