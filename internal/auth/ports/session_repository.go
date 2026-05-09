package ports

import (
	"context"
	"time"
)

type SessionRepository interface {
	Set(ctx context.Context, key string, val string, ttl time.Duration) (status string, err error)
	Get(ctx context.Context, key string) (val string, err error)
	Delete(ctx context.Context, key string) error
}
