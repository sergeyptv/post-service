package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Username string        `env:"USERNAME" env-prefix:"REDIS_" env-required`
	Password string        `env:"PASSWORD" env-prefix:"REDIS_" env-required`
	Host     string        `env:"HOST" env-prefix:"REDIS_" env-required`
	Port     string        `env:"PORT" env-prefix:"REDIS_" env-required`
	DB       int           `env:"DB" env-prefix:"REDIS_" env-required`
	TokenTtl time.Duration `env:"TOKEN_TTL" env-prefix:"REDIS_" env-required`

	IpRateLimit         int           `env:"IP_RATE_LIMIT" env-prefix:"REDIS_" env-required`
	IpRateLimiterTtl    time.Duration `env:"IP_RATE_LIMITER_TTL" env-prefix:"REDIS_" env-required`
	EmailRateLimit      int           `env:"EMAIL_RATE_LIMIT" env-prefix:"REDIS_" env-required`
	EmailRateLimiterTtl time.Duration `env:"EMAIL_RATE_LIMITER_TTL" env-prefix:"REDIS_" env-required`
}

type Client struct {
	Db *redis.Client
}

func New(ctx context.Context, c Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,
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
