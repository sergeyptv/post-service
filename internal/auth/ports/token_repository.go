package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, token domain.Token) (string, error)
	GetToken(ctx context.Context, tokenUuid string) (domain.Token, error)
	UpdateToken(ctx context.Context, tokenUuid string, updToken domain.UpdateToken) error
	DeleteToken(ctx context.Context, tokenUuid string) error
}
