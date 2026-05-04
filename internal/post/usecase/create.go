package usecase

import (
	"context"
	"fmt"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"github.com/sergeyptv/post_service/internal/post/domain"
	"log/slog"
)

func (p *post) Create(ctx context.Context, user domain.User, post domain.Post) (string, error) {
	const op = "usecase.Create"

	log := p.log.With(slog.String("op", op))

	postUuid, err := p.postRepository.Create(ctx, user, post)
	if err != nil {
		log.Error("Failed to create new post", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return postUuid, nil
}
