package main

import (
	"context"
	"github.com/sergeyptv/post_service/internal/outbox/config"
	"github.com/sergeyptv/post_service/internal/outbox/producer/kafka"
	"github.com/sergeyptv/post_service/internal/outbox/repository/postgres"
	"github.com/sergeyptv/post_service/internal/outbox/usecase"
	"github.com/sergeyptv/post_service/internal/platform/kafka_produce"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	platformPostgres "github.com/sergeyptv/post_service/internal/platform/postgres"
	"github.com/sergeyptv/post_service/internal/platform/transaction"
	"log/slog"
	"math"
	"math/rand"
	"os/signal"
	"syscall"
	"time"
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

	postgresOutboxRepository := postgres.NewPostgresOutboxRepository(pool.Db)

	p, err := kafka_produce.NewProducer(ctx, cfg.KafkaProducer)
	if err != nil {
		return err
	}
	defer p.Close()

	kafkaEventProducer := kafka.NewKafkaEventProducer(p)

	txWrapper := transaction.New(pool.Db)

	outboxService := usecase.NewOutboxService(log, cfg, postgresOutboxRepository, kafkaEventProducer, txWrapper)

	backoff := time.Second
	const maxBackoff = 30 * time.Second

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = outboxService.Worker(ctx)
		if err != nil {
			log.Error("Worker failed", logger.Error(err))

			time.Sleep(backoff)

			backoff = time.Duration(math.Min(float64(backoff*2), float64(maxBackoff)))
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		backoff = time.Second

		sleep := time.Duration(cfg.WorkerFrequencySec)*time.Second + time.Duration(rand.Intn(500))*time.Millisecond

		time.Sleep(sleep)
	}
}
