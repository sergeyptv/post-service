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

func (p *post) Update(ctx context.Context, post domain.Post) error {
	const op = "usecase.Update"

	log := p.log.With(slog.String("op", op))

	err := p.postRepository.Update(ctx, post)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Failed to update unknown post", logger.Error(err))

			return fmt.Errorf("%s: %w", op, domain.ErrPostNotExist)
		}

		log.Error("Failed to update post", logger.Error(err))

		return fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return nil
}
