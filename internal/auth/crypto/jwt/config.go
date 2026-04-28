package jwt

import "time"

type Config struct {
	PrivateKey string        `env:"PRIVATE_KEY" env-prefix:"TOKEN_" env-required`
	PublicKey  string        `env:"PUBLIC_KEY" env-prefix:"TOKEN_" env-required`
	TokenTTL   time.Duration `env:"TTL" env-prefix:"TOKEN_" env-default:"1h"`
}
