package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sergeyptv/post_service/auth/internal/repository"
	platformRedis "github.com/sergeyptv/post_service/platform/redis"
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

var refreshScript = redis.NewScript(`
	local old = redis.call("GET", KEYS[1])
	if not old then
		return 0
	end
	
	redis.call("DEL", KEYS[1])
	redis.call("SET", KEYS[2], ARGV[1], "PX", ARGV[2])
	
	return 1
`)

func (r *redisSessionRepository) SetToken(ctx context.Context, jti string, refreshToken string, ttl time.Duration) (status string, err error) {
	const op = "repository.redis.SetToken"

	status, err = r.client.Db.Set(ctx, jti, refreshToken, ttl).Result()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (r *redisSessionRepository) RotateToken(ctx context.Context, oldJti, newJti, refreshToken string, ttl time.Duration) (success bool, err error) {
	const op = "repository.redis.RotateToken"

	res, err := refreshScript.Run(ctx, r.client.Db, []string{oldJti, newJti}, refreshToken, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	resInt64, ok := res.(int64)
	if !ok {
		return false, fmt.Errorf("%s: %w", op, repository.ErrGetResult)
	}

	return resInt64 == 1, nil
}

func (r *redisSessionRepository) DeleteToken(ctx context.Context, jti string) error {
	const op = "repository.redis.DeleteToken"

	err := r.client.Db.Del(ctx, jti).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
