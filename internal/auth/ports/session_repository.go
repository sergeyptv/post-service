package ports

import (
	"context"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, userUuid string, token string) (string, error)
	GetToken(ctx context.Context, userUuid string) (string, error)
	UpdateToken(ctx context.Context, jti string, newToken string) error
}
