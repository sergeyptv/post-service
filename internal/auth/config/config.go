package config

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
	authJwt "github.com/sergeyptv/post_service/internal/auth/crypto/jwt"
	"github.com/sergeyptv/post_service/internal/platform/config"
	"github.com/sergeyptv/post_service/internal/platform/grpcServer"
	"github.com/sergeyptv/post_service/internal/platform/httpserver"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
	"github.com/sergeyptv/post_service/internal/platform/redis"
	"os"
)

type Config struct {
	App        config.App
	Jwt        authJwt.Config
	Postgres   postgres.Config
	Redis      redis.Config
	HttpServer httpserver.Config
	GrpcServer grpcServer.Config
}

func MustLoad() *Config {
	return mustParseEnv()
}

func mustParseEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot get all env")
	}

	cfg.Jwt.PrivateKey, cfg.Jwt.PublicKey = mustParseRSAKeys(cfg.Jwt)

	return &cfg
}

func mustParseRSAKeys(c authJwt.Config) (*rsa.PrivateKey, *rsa.PublicKey) {
	filePrivateKey, err := os.OpenFile(fmt.Sprintf("../../%s", c.PrivateKeyPath), os.O_RDONLY, 0777)
	if err != nil {
		panic("cannot open rsa private key")
	}

	var rsaPrivateKeyBytes []byte
	n, err := filePrivateKey.Read(rsaPrivateKeyBytes)
	if err != nil || n == 0 {
		panic("cannot read rsa private key")
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKeyBytes)
	if err != nil {
		panic("cannot parse rsa private key")
	}

	filePublicKey, err := os.OpenFile(fmt.Sprintf("../../%s", c.PublicKeyPath), os.O_RDONLY, 0777)
	if err != nil {
		panic("cannot open rsa public key")
	}

	var rsaPublicKeyBytes []byte
	n, err = filePublicKey.Read(rsaPublicKeyBytes)
	if err != nil || n == 0 {
		panic("cannot read rsa public key")
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKeyBytes)
	if err != nil {
		panic("cannot parse rsa public key")
	}

	return rsaPrivateKey, rsaPublicKey
}
