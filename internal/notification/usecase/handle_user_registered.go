package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/notification/repository"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"log/slog"
)

func (n *notification) HandleUserRegistered(ctx context.Context) error {
	const op = "usecase.HandleUserRegistered"

	log := n.log.With(slog.String("op", op))

	commit := func() error {
		err := n.poller.Commit()
		if err != nil {
			log.Error("Failed to commit getting event", logger.Error(err))

			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}

	staleEvent, err := n.notificationRepository.GetStaleEvent(ctx)
	if err != nil {
		log.Error("Failed to get a stale event from db", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	if staleEvent.Uuid != "" {
		log.Info("Got new event", slog.Attr{
			Key:   "user.registered",
			Value: slog.StringValue(fmt.Sprintf("%+v", staleEvent)),
		})

		err = n.notificationRepository.MarkSuccess(ctx, staleEvent.Uuid)
		if err != nil {
			log.Error("Failed to set success status to the event", logger.Error(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	event, err := n.poller.Poll()
	if err != nil {
		log.Error("Failed to poll event", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = n.notificationRepository.TryProcess(ctx, event)
	if err != nil {
		if errors.Is(err, repository.ErrEventAlreadySuccess) ||
			errors.Is(err, repository.ErrEventStatusProcessing) {
			return commit()
		}
		log.Error("Failed to try process the event", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Got new event", slog.Attr{
		Key:   "user.registered",
		Value: slog.StringValue(fmt.Sprintf("%+v", event)),
	})

	err = n.notificationRepository.MarkSuccess(ctx, event.Uuid)
	if err != nil {
		log.Error("Failed to set success status to the event", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return commit()
}
