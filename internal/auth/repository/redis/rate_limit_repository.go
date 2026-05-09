package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	platformRedis "github.com/sergeyptv/post_service/internal/platform/redis"
	"time"
)

type redisRateLimitRepository struct {
	client *platformRedis.Client
}

func NewRedisRateLimitRepository(client *platformRedis.Client) *redisRateLimitRepository {
	return &redisRateLimitRepository{
		client: client,
	}
}

func (p *redisRateLimitRepository) SetLimit(ctx context.Context, limiter string, limit string, ttl time.Duration) (status string, err error) {
	const op = "repository.redis.SetLimit"

	status, err = p.client.Db.Set(ctx, limiter, limit, ttl).Result()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (p *redisRateLimitRepository) GetLimit(ctx context.Context, limiter string) (limit string, err error) {
	const op = "repository.redis.GetLimit"

	err = p.client.Db.Get(ctx, limiter).Scan(&limit)
	if err != nil {
		if errors.Is(err, redis.ErrClosed) {
			return "", fmt.Errorf("%s: %w", op, repository.ErrDbClientClosed)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return limit, nil
}
