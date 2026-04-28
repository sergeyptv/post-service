package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/post/domain"
)

type PostRepository interface {
	Create(ctx context.Context, post domain.Post) (string, error)
	Get(ctx context.Context, postUuid string) (domain.Post, error)
	List(ctx context.Context, username string) ([]string, error)
	Update(ctx context.Context, post domain.Post) error
	Delete(ctx context.Context, postUuid string) error
}
