package postgres

import (
	"context"
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

func (p *postgresTokenRepository) CreateToken(ctx context.Context, token domain.Token) (string, error) {
}
func (p *postgresTokenRepository) GetToken(ctx context.Context, tokenUuid string) (domain.Token, error) {
}
func (p *postgresTokenRepository) UpdateToken(ctx context.Context, tokenUuid string, updToken domain.UpdateToken) error {
}
func (p *postgresTokenRepository) DeleteToken(ctx context.Context, tokenUuid string) error {}
