package ports

import (
	"context"
	"time"
)

type SessionRepository interface {
	SetToken(ctx context.Context, jti string, refreshToken string, ttl time.Duration) (status string, err error)
	DeleteToken(ctx context.Context, jti string) error
	RotateToken(ctx context.Context, oldJti, newJti, refreshToken string, ttl time.Duration) (success bool, err error)
}
