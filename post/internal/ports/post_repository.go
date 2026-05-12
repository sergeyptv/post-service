package ports

import (
	"context"
	"github.com/sergeyptv/post_service/post/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, userUuid, username string, post domain.Post) (string, error)
	Get(ctx context.Context, userUuid string, postUuid string) (domain.Post, error)
	List(ctx context.Context, userUuid string) ([]string, error)
	Update(ctx context.Context, userUuid string, post domain.Post) error
	Delete(ctx context.Context, userUuid, postUuid string) error
}
