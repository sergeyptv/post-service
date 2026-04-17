package jwt

import "time"

type Config struct {
	TokenKey string        `env:"KEY" env-prefix:"TOKEN_" env-required`
	TokenTTL time.Duration `env:"TTL" env-prefix:"TOKEN_" env-default:"1h"`
}
