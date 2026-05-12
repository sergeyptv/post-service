package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/platform/logger"
	"github.com/sergeyptv/post_service/post/internal/domain"
	"log/slog"
)

func (p *post) Update(ctx context.Context, user domain.User, post domain.Post) error {
	const op = "usecase.Update"

	log := p.log.With(slog.String("op", op))

	err := p.postRepository.Update(ctx, user.Uuid, post)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("Failed to update unknown post", logger.Error(err))

			return fmt.Errorf("%s: %w", op, domain.ErrPostNotExist)
		}

		log.Error("Failed to update post", logger.Error(err))

		return fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return nil
}
