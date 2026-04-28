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

func (p *post) Get(ctx context.Context, postUuid string) (domain.Post, error) {
	const op = "usecase.Get"

	log := p.log.With(slog.String("op", op))

	post, err := p.postRepository.Get(ctx, postUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Failed to get unknown post", logger.Error(err))

			return domain.Post{}, fmt.Errorf("%s: %w", op, domain.ErrPostNotExist)
		}

		log.Error("Failed to get post", logger.Error(err))

		return domain.Post{}, fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return post, nil
}
