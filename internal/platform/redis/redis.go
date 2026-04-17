package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Username string `env:"USERNAME" env-prefix:"REDIS_" env-required`
	Password string `env:"PASSWORD" env-prefix:"REDIS_" env-required`
	Host     string `env:"HOST" env-prefix:"REDIS_" env-required`
	Port     string `env:"PORT" env-prefix:"REDIS_" env-required`
	DB       int    `env:"DB" env-prefix:"REDIS_" env-required`
}

type Client struct {
	db *redis.Client
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
		db: rdb,
	}, nil
}

func (r *Client) Close() {
	if err := r.db.Close(); err != nil {
		fmt.Printf("Error close redis: %s\n", err)
	}
}
