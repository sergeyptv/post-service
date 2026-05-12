package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sergeyptv/post_service/auth/internal/repository"
	platformRedis "github.com/sergeyptv/post_service/platform/redis"
	"time"
)

var rateLimitScript = redis.NewScript(`
	local current = redis.call("INCR", KEYS[1])
	if current == 1 then
		redis.call("PEXPIRE", KEYS[1], ARGV[1])
	end
	return current
`)

type redisRateLimitRepository struct {
	client *platformRedis.Client
}

func NewRedisRateLimitRepository(client *platformRedis.Client) *redisRateLimitRepository {
	return &redisRateLimitRepository{
		client: client,
	}
}

func (p *redisRateLimitRepository) Increment(ctx context.Context, limiter string, ttl time.Duration) (limit int64, err error) {
	const op = "repository.redis.Increment"

	res, err := rateLimitScript.Run(ctx, p.client.Db, []string{limiter}, ttl.Milliseconds()).Result()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	limit, ok := res.(int64)
	if !ok {
		return 0, fmt.Errorf("%s: %w", op, repository.ErrGetResult)
	}

	return limit, nil
}
