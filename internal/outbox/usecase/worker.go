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

	publishedEventsMap := make(map[string]bool)

	for _, event := range events {
		err = o.publisher.Publish(event)
		if err != nil {
			publishedEventsMap[event.Uuid] = false
		} else {
			publishedEventsMap[event.Uuid] = true
		}
	}

	for eventUuid, published := range publishedEventsMap {
		if published {
			terr := o.outboxRepository.MarkSent(ctx, eventUuid)
			if terr != nil {
				log.Error("Failed to mark events sent", logger.Error(terr))

				return fmt.Errorf("%s: %w", op, terr)
			}
		} else {
			terr := o.outboxRepository.MarkFailed(ctx, eventUuid)
			if terr != nil {
				log.Error("Failed to mark events failed", logger.Error(terr))

				return fmt.Errorf("%s: %w", op, terr)
			}
		}
	}

	return nil
}
