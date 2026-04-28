package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"github.com/sergeyptv/post_service/internal/post/domain"
	"log/slog"
)

func (p *post) Delete(ctx context.Context, postUuid string) error {
	const op = "usecase.Delete"

	log := p.log.With(slog.String("op", op))

	err := p.postRepository.Delete(ctx, postUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Failed to delete unknown post", logger.Error(err))

			return fmt.Errorf("%s: %w", op, domain.ErrPostNotExist)
		}

		log.Error("Failed to delete post", logger.Error(err))

		return fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return nil
}
