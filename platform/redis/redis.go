package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Username string        `env:"USERNAME"  env-required`
	Password string        `env:"PASSWORD" env-required`
	Host     string        `env:"HOST" env-required`
	Port     string        `env:"PORT" env-required`
	DB       int           `env:"DB" env-required`
	TokenTtl time.Duration `env:"TOKEN_TTL" env-required`

	IpRateLimit         int64         `env:"IP_RATE_LIMIT" env-required`
	IpRateLimiterTtl    time.Duration `env:"IP_RATE_LIMITER_TTL" env-required`
	EmailRateLimit      int64         `env:"EMAIL_RATE_LIMIT" env-required`
	EmailRateLimiterTtl time.Duration `env:"EMAIL_RATE_LIMITER_TTL" env-required`
}

type Client struct {
	Db *redis.Client
}

func New(ctx context.Context, c Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", c.Host, c.Port),
		DB:   c.DB,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &Client{
		Db: rdb,
	}, nil
}

func (r *Client) Close() {
	if err := r.Db.Close(); err != nil {
		fmt.Printf("Error close redis: %s\n", err)
	}
}
