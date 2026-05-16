package cache

import "time"

type Config struct {
	Ttl time.Duration `env:"TTL" env-required`
}
