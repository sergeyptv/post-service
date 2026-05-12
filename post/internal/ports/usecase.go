package ports

import (
	"context"
	"github.com/sergeyptv/post_service/post/internal/domain"
)

type Usecase interface {
	Create(ctx context.Context, user domain.User, post domain.Post) (string, error)
	Get(ctx context.Context, user domain.User, postUuid string) (domain.Post, error)
	List(ctx context.Context, user domain.User) ([]string, error)
	Update(ctx context.Context, user domain.User, post domain.Post) error
	Delete(ctx context.Context, user domain.User, postUuid string) error
}
