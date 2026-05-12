package ports

import (
	"context"
	"time"
)

type RateLimitRepository interface {
	Increment(ctx context.Context, limiter string, ttl time.Duration) (limit int64, err error)
}
