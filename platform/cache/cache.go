package cache

import "time"

type Config struct {
	Ttl time.Duration `env:"TTL" env-prefix:"CACHE_" env-required`
}
