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

	var events []domain.UserRegisteredEvent

	err := o.txWrapper.Wrap(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var terr error

		events, terr = o.outboxRepository.GetPending(ctx, tx, o.cfg.BatchSize)
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

	publishEventsSuccess := make([]string, 0, len(events))
	publishEventsFailed := make([]string, 0, len(events))

	for _, event := range events {
		err = o.publisher.Publish(event)
		if err != nil {
			publishEventsFailed = append(publishEventsFailed, event.Uuid)
		} else {
			publishEventsSuccess = append(publishEventsSuccess, event.Uuid)
		}
	}

	if len(publishEventsSuccess) > 0 {
		err = o.outboxRepository.MarkSent(ctx, publishEventsSuccess)
		if err != nil {
			log.Error("Failed to mark events sent", logger.Error(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if len(publishEventsFailed) > 0 {
		err = o.outboxRepository.MarkFailed(ctx, publishEventsFailed)
		if err != nil {
			log.Error("Failed to mark events failed", logger.Error(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
