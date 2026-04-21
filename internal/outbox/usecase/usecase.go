package usecase

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/internal/outbox/ports"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"log/slog"
)

type outbox struct {
	log              *slog.Logger
	outboxRepository ports.OutboxRepository
	publisher        ports.EventPublisher
	txWrapper        ports.TransactionWrapper
}

func NewOutboxService(log *slog.Logger, outboxRepository ports.OutboxRepository, publisher ports.EventPublisher, txWrapper ports.TransactionWrapper) *outbox {
	return &outbox{
		log:              log,
		outboxRepository: outboxRepository,
		publisher:        publisher,
		txWrapper:        txWrapper,
	}
}

func (o *outbox) Worker(ctx context.Context) error {
	const op = "usecase.Worker"

	log := o.log.With(slog.String("op", op))

	err := o.txWrapper.Wrap(ctx, func(ctx context.Context, tx pgx.Tx) error {
		events, err := o.outboxRepository.GetPending(ctx, tx, 20)
		if err != nil {
			log.Error("Failed to get pending", logger.Error(err))

			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	userUuids := make([]string, len(events))
	for _, event := range events {
		userUuids = append(userUuids, event.UserUuid)
	}

	err = o.publisher.Publish(events)
	if err != nil {
		err = o.txWrapper.Wrap(ctx, func(ctx context.Context, tx pgx.Tx) error {
			err = o.outboxRepository.MarkFailed(ctx, tx, userUuids)
			if err != nil {
				log.Error("Failed to mark events failed", logger.Error(err))

				return fmt.Errorf("%s: %w", op, err)
			}

			return nil
		})
		log.Error("Failed to publish events", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	} else {
		err = o.txWrapper.Wrap(ctx, func(ctx context.Context, tx pgx.Tx) error {
			err = o.outboxRepository.MarkSent(ctx, tx, userUuids)
			if err != nil {
				log.Error("Failed to mark events sent", logger.Error(err))

				return fmt.Errorf("%s: %w", op, err)
			}

			return nil
		})
	}

	return nil
}
