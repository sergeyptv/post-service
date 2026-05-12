package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/platform/logger"
	"github.com/sergeyptv/post_service/post/internal/domain"
	"log/slog"
)

func (p *post) Get(ctx context.Context, user domain.User, postUuid string) (domain.Post, error) {
	const op = "usecase.Get"

	log := p.log.With(slog.String("op", op))

	post, err := p.postRepository.Get(ctx, user.Uuid, postUuid)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotExist) {
			log.Warn("Failed to get unknown post", logger.Error(err))

			return domain.Post{}, fmt.Errorf("%s: %w", op, domain.ErrPostNotExist)
		}

		log.Error("Failed to get post", logger.Error(err))

		return domain.Post{}, fmt.Errorf("%s: %w", op, domain.ErrBadGateway)
	}

	return post, nil
}
