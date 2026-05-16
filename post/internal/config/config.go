package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sergeyptv/post_service/platform/cache"
	"github.com/sergeyptv/post_service/platform/config"
	grpcClient "github.com/sergeyptv/post_service/platform/grpc_client"
	httpServer "github.com/sergeyptv/post_service/platform/http_server"
	"github.com/sergeyptv/post_service/platform/jwt"
	"github.com/sergeyptv/post_service/platform/postgres"
)

type Config struct {
	App        config.App        `env-prefix:"APP_"`
	Cache      cache.Config      `env-prefix:"CACHE_"`
	Server     httpServer.Config `env-prefix:"HTTP_"`
	GrpcClient grpcClient.Config `env-prefix:"GRPC_CLIENT_"`
	Postgres   postgres.Config   `env-prefix:"POSTGRES_"`
	Jwt        jwt.ConfigParser  `env-prefix:"TOKEN_"`
}

func MustLoad() *Config {
	return mustParseEnv()
}

func mustParseEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
