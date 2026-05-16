package config

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sergeyptv/post_service/platform/config"
	grpcServer "github.com/sergeyptv/post_service/platform/grpc_server"
	httpServer "github.com/sergeyptv/post_service/platform/http_server"
	authJwt "github.com/sergeyptv/post_service/platform/jwt"
	"github.com/sergeyptv/post_service/platform/postgres"
	"github.com/sergeyptv/post_service/platform/redis"
)

type Config struct {
	App        config.App           `env-prefix:"APP_"`
	Jwt        authJwt.ConfigSigner `env-prefix:"TOKEN_"`
	Postgres   postgres.Config      `env-prefix:"POSTGRES_"`
	Redis      redis.Config         `env-prefix:"REDIS_"`
	HttpServer httpServer.Config    `env-prefix:"HTTP_"`
	GrpcServer grpcServer.Config    `env-prefix:"GRPC_SERVER_"`
}

func MustLoad() *Config {
	return mustParseEnv()
}

func mustParseEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	cfg.Jwt.PrivateKey, cfg.Jwt.PublicKey = mustParseRSAKeys(cfg.Jwt)

	return &cfg
}

func mustParseRSAKeys(c authJwt.ConfigSigner) (*rsa.PrivateKey, *rsa.PublicKey) {
	var rsaPrivateKeyBytes []byte

	rsaPrivateKeyBytes, err := os.ReadFile(c.PrivateKeyPath)
	if err != nil || len(rsaPrivateKeyBytes) == 0 {
		panic(fmt.Sprintf("cannot read rsa private key %s\n", err))
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKeyBytes)
	if err != nil {
		panic("cannot parse rsa private key")
	}

	var rsaPublicKeyBytes []byte

	rsaPublicKeyBytes, err = os.ReadFile(c.PublicKeyPath)
	if err != nil || len(rsaPublicKeyBytes) == 0 {
		panic("cannot read rsa public key")
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKeyBytes)
	if err != nil {
		panic("cannot parse rsa public key")
	}

	return rsaPrivateKey, rsaPublicKey
}
