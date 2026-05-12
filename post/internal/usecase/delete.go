package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/platform/logger"
	"github.com/sergeyptv/post_service/post/internal/domain"
	"log/slog"
)

func (p *post) Delete(ctx context.Context, user domain.User, postUuid string) error {
	const op = "usecase.Delete"

	log := p.log.With(slog.String("op", op))

	err := p.postRepository.Delete(ctx, user.Uuid, postUuid)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotExist) {
			log.Warn("Failed to delete unknown post", logger.Error(err))

			return fmt.Errorf("%s: %w", op, domain.ErrPostNotExist)
		}

		log.Error("Failed to delete post", logger.Error(err))

		return fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return nil
}
