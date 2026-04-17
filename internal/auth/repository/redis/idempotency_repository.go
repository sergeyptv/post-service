package redis

import (
	"github.com/sergeyptv/post_service/internal/platform/redis"
)

type redisIdempotencyRepositry struct {
	client *redis.Client
}

func NewRedisIdempotencyRepositry(client *redis.Client) *redisIdempotencyRepositry {
	return &redisIdempotencyRepositry{
		client: client,
	}
}
