package usecase

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/internal/outbox/domain"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"log/slog"
)

func (o *outbox) Worker(ctx context.Context) error {
	const op = "usecase.Worker"

	log := o.log.With(slog.String("op", op))

	var eventUuids []string
	var events []domain.UserRegisteredEvent

	err := o.txWrapper.Wrap(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var terr error

		eventUuids, events, terr = o.outboxRepository.GetPending(ctx, tx, 20)
		if terr != nil {
			return terr
		}

		return nil
	})
	if err != nil {
		log.Error("Failed to get pending", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	if len(events) == 0 {
		return nil
	}

	err = o.publisher.Publish(events)
	if err != nil {
		terr := o.outboxRepository.MarkFailed(ctx, eventUuids)
		if terr != nil {
			log.Error("Failed to mark events failed", logger.Error(terr))

			return fmt.Errorf("%s: %w", op, terr)
		}

		log.Error("Failed to publish events", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	terr := o.outboxRepository.MarkSent(ctx, eventUuids)
	if terr != nil {
		log.Error("Failed to mark events sent", logger.Error(terr))

		return fmt.Errorf("%s: %w", op, terr)
	}

	return nil
}
