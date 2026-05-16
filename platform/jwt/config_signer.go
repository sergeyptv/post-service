package jwt

import (
	"crypto/rsa"
	"time"
)

type ConfigSigner struct {
	PrivateKeyPath  string `env:"PRIVATE_KEY_PATH" env-required`
	PrivateKey      *rsa.PrivateKey
	PublicKeyPath   string `env:"PUBLIC_KEY_PATH" env-required`
	PublicKey       *rsa.PublicKey
	Issuer          string        `env:"ISSUER" env-required`
	Format          string        `env:"FORMAT" env-required`
	Algorithm       string        `env:"ALGORITHM" env-required`
	Kid             string        `env:"KID" env-required`
	AccessTokenTtl  time.Duration `env:"ACCESS_TTL" env-required`
	RefreshTokenTtl time.Duration `env:"REFRESH_TTL" env-required`
}
