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

func (p *post) List(ctx context.Context, username string) ([]string, error) {
	const op = "usecase.List"

	log := p.log.With(slog.String("op", op))

	postUuids, err := p.postRepository.List(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Failed to list posts for specified username", logger.Error(err))

			return []string{}, fmt.Errorf("%s: %w", op, domain.ErrPostNotExist)
		}

		log.Error("Failed to list posts", logger.Error(err))

		return []string{}, fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return postUuids, nil
}
