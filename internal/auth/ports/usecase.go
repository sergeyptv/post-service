package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

type Usecase interface {
	Register(ctx context.Context, user domain.User) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}
