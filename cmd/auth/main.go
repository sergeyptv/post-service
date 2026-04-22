package main

import (
	"context"
	"github.com/sergeyptv/post_service/internal/auth/config"
	"github.com/sergeyptv/post_service/internal/auth/crypto/jwt"
	"github.com/sergeyptv/post_service/internal/auth/repository/postgres"
	"github.com/sergeyptv/post_service/internal/auth/repository/redis"
	"github.com/sergeyptv/post_service/internal/auth/usecase"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	platformPostgres "github.com/sergeyptv/post_service/internal/platform/postgres"
	platformRedis "github.com/sergeyptv/post_service/internal/platform/redis"
	"github.com/sergeyptv/post_service/internal/platform/transaction"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.App.Env)

	if err := appRun(log, cfg); err != nil {
		panic(err)
	}
}

func appRun(log *slog.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jwtTokenSigner := jwt.NewJwtTokenSigner(cfg.Jwt)

	pool, err := platformPostgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pool.Close()

	postgresUserRepository := postgres.NewPostgresUserRepository(pool)
	postgresTokenRepository := postgres.NewPostgresTokenRepository(pool)
	postgresOutboxRepository := postgres.NewPostgresOutboxRepository(pool)

	client, err := platformRedis.New(ctx, cfg.Redis)
	if err != nil {
		return err
	}
	defer client.Close()

	redisIdempotencyRepository := redis.NewRedisIdempotencyRepositry(client)

	txWrapper := transaction.New(pool.Db)

	authService := usecase.NewAuthService(log, postgresUserRepository, postgresOutboxRepository, postgresTokenRepository, redisIdempotencyRepository, jwtTokenSigner, txWrapper)
}
