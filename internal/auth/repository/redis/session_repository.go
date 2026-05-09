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

func (r *redisSessionRepository) SetToken(ctx context.Context, jti string, refreshToken string, ttl time.Duration) (status string, err error) {
	const op = "repository.redis.SetToken"

	status, err = r.client.Db.Set(ctx, jti, refreshToken, ttl).Result()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (r *redisSessionRepository) GetToken(ctx context.Context, jti string) (refreshToken string, err error) {
	const op = "repository.redis.GetToken"

	err = r.client.Db.Get(ctx, jti).Scan(&refreshToken)
	if err != nil {
		if errors.Is(err, redis.ErrClosed) {
			return "", fmt.Errorf("%s: %w", op, repository.ErrDbClientClosed)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, nil
}

func (r *redisSessionRepository) DeleteToken(ctx context.Context, jti string) error {
	const op = "repository.redis.DeleteToken"

	err := r.client.Db.Del(ctx, jti).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
