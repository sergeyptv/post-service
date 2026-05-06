package jwt

import (
	"crypto/rsa"
	"time"
)

type ConfigSigner struct {
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH" env-prefix:"TOKEN_" env-required`
	PrivateKey     *rsa.PrivateKey
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH" env-prefix:"TOKEN_" env-required`
	PublicKey      *rsa.PublicKey
	Issuer         string        `env:"ISSUER" env-prefix:"TOKEN_" env-required`
	Format         string        `env:"FORMAT" env-prefix:"TOKEN_" env-required`
	Algorithm      string        `env:"ALGORITHM" env-prefix:"TOKEN_" env-required`
	TokenTTL       time.Duration `env:"TTL" env-prefix:"TOKEN_" env-default:"1h"`
}
