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

func (p *post) List(ctx context.Context, user domain.User) ([]string, error) {
	const op = "usecase.List"

	log := p.log.With(slog.String("op", op))

	postUuids, err := p.postRepository.List(ctx, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}

		log.Error("Failed to list posts", logger.Error(err))

		return []string{}, fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return postUuids, nil
}
