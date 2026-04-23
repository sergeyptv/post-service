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

	event, err := n.poller.Poll()
	if err != nil {
		log.Error("Failed to poll event", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	commit := func() error {
		err = n.poller.Commit()
		if err != nil {
			log.Error("Failed to commit getting event", logger.Error(err))

			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}

	err = n.notificationRepository.TryProcess(ctx, event.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrEventAlreadySuccess) {
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
