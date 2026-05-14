package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/sergeyptv/post_service/notification/internal/config"
	"github.com/sergeyptv/post_service/notification/internal/consumer/kafka"
	"github.com/sergeyptv/post_service/notification/internal/repository/postgres"
	"github.com/sergeyptv/post_service/notification/internal/usecase"
	"github.com/sergeyptv/post_service/platform/kafka_consume"
	"github.com/sergeyptv/post_service/platform/logger"
	platformPostgres "github.com/sergeyptv/post_service/platform/postgres"
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

	pool, err := platformPostgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer pool.Close()

	postgresNotificationRepository := postgres.NewPostgresNotificationRepository(pool.Db)

	c, err := kafka_consume.NewConsumer(ctx, cfg.KafkaConsumer)
	if err != nil {
		return err
	}
	defer c.Close()

	kafkaEventConsumer := kafka.NewKafkaEventConsumer(c)

	notificationService := usecase.NewNotificationService(log, cfg, postgresNotificationRepository, kafkaEventConsumer)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
		}

		err = notificationService.HandleUserRegistered(ctx)
		if err != nil {
			log.Error("User registered handler failed", logger.Error(err))

			continue
		}
	}
}
