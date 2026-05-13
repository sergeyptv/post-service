package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sergeyptv/post_service/platform/cache"
	"github.com/sergeyptv/post_service/platform/config"
	grpcClient "github.com/sergeyptv/post_service/platform/grpc_client"
	httpServer "github.com/sergeyptv/post_service/platform/http_server"
	"github.com/sergeyptv/post_service/platform/jwt"
	"github.com/sergeyptv/post_service/platform/migrator"
	"github.com/sergeyptv/post_service/platform/postgres"
)

type Config struct {
	App        config.App
	Migrator   migrator.Config
	Cache      cache.Config
	Server     httpServer.Config
	GrpcClient grpcClient.Config
	Postgres   postgres.Config
	Jwt        jwt.ConfigParser
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
