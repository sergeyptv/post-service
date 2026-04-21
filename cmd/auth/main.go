package main

import (
	"context"
	"github.com/sergeyptv/post_service/internal/auth/config"
	"github.com/sergeyptv/post_service/internal/auth/crypto/jwt"
	"github.com/sergeyptv/post_service/internal/auth/repository/postgres"
	"github.com/sergeyptv/post_service/internal/auth/repository/redis"
	"github.com/sergeyptv/post_service/internal/auth/usecase"
	"github.com/sergeyptv/post_service/internal/outbox/producer/kafka"
	"github.com/sergeyptv/post_service/internal/platform/kafka_produce"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	pPostgres "github.com/sergeyptv/post_service/internal/platform/postgres"
	pRedis "github.com/sergeyptv/post_service/internal/platform/redis"
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

	p, deliveryStatus, err := kafka_produce.NewProducer(ctx, cfg.KafkaProducer)
	if err != nil {
		return err
	}
	defer p.Close()

	go func(ctx context.Context, deliveryStatus chan string) {
		for {
			select {
			case s := <-deliveryStatus:
				log.Info("kafka delivery status: %s\n", s)
			case <-ctx.Done():
				log.Info("Program context exceeded: %s\n", ctx.Err())
			}
		}
	}(ctx, deliveryStatus)

	kafkaEventProducer := kafka.NewKafkaEventProducer(p)

	pool, err := pPostgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pool.Close()

	postgresUserRepository := postgres.NewPostgresUserRepository(pool)
	postgresTokenRepository := postgres.NewPostgresTokenRepository(pool)

	client, err := pRedis.New(ctx, cfg.Redis)
	if err != nil {
		return err
	}
	defer client.Close()

	redisIdempotencyRepository := redis.NewRedisIdempotencyRepositry(client)

	authService := usecase.NewAuthService(log, postgresUserRepository, postgresTokenRepository, kafkaEventProducer, redisIdempotencyRepository, jwtTokenSigner)
}
