package main

import (
	"context"
	"fmt"
	"github.com/sergeyptv/post_service/notification/internal/config"
	"github.com/sergeyptv/post_service/notification/internal/consumer/kafka"
	"github.com/sergeyptv/post_service/notification/internal/repository/postgres"
	"github.com/sergeyptv/post_service/notification/internal/usecase"
	"github.com/sergeyptv/post_service/platform/kafka_consume"
	"github.com/sergeyptv/post_service/platform/logger"
	"github.com/sergeyptv/post_service/platform/migrator"
	platformPostgres "github.com/sergeyptv/post_service/platform/postgres"
	"log/slog"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	dsn := fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	err := migrator.Up(cfg.Migrator.Dir, dsn)
	if err != nil {
		panic(err)
	}

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
