package jwt

import (
	"crypto/rsa"
	"time"
)

type Config struct {
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH" env-prefix:"TOKEN_" env-required`
	PrivateKey     *rsa.PrivateKey
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH" env-prefix:"TOKEN_" env-required`
	PublicKey      *rsa.PublicKey
	TokenTTL       time.Duration `env:"TTL" env-prefix:"TOKEN_" env-default:"1h"`
}
