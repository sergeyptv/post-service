package ports

import (
	"context"
	"time"
)

type RateLimitRepository interface {
	SetLimit(ctx context.Context, limiter string, limit string, ttl time.Duration) (status string, err error)
	GetLimit(ctx context.Context, limiter string) (limit string, err error)
}
