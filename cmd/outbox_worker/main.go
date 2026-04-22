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
	"math/rand"
	"os"
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

				close(deliveryStatus)
				return
			}
		}
	}(ctx, deliveryStatus)

	kafkaEventProducer := kafka.NewKafkaEventProducer(p)

	txWrapper := transaction.New(pool.Db)

	outboxService := usecase.NewOutboxService(log, postgresOutboxRepository, kafkaEventProducer, txWrapper)

	backoff := time.Second
	const maxBackoff = 30 * time.Second

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

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		backoff = time.Second

		time.Sleep(time.Duration(cfg.WorkerFrequencySec) * time.Second)

		sleep := backoff + time.Duration(rand.Intn(500))*time.Millisecond
		time.Sleep(sleep)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	return nil
}
