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

type redisSessionRepository struct {
	client *platformRedis.Client
}

func NewRedisSessionRepository(client *platformRedis.Client) *redisSessionRepository {
	return &redisSessionRepository{
		client: client,
	}
}

func (p *redisSessionRepository) Set(ctx context.Context, key string, val string, ttl time.Duration) (status string, err error) {
	const op = "repository.redis.Set"

	status, err = p.client.Db.Set(ctx, key, val, ttl).Result()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (p *redisSessionRepository) Get(ctx context.Context, key string) (val string, err error) {
	const op = "repository.redis.Get"

	err = p.client.Db.Get(ctx, key).Scan(&val)
	if err != nil {
		if errors.Is(err, redis.ErrClosed) {
			return "", fmt.Errorf("%s: %w", op, repository.ErrDbClientClosed)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return val, nil
}

func (p *redisSessionRepository) Delete(ctx context.Context, key string) error {
	const op = "repository.redis.Delete"

	err := p.client.Db.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
