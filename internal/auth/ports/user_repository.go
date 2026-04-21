package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.CreateUser) (string, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateEvent(ctx context.Context, event domain.UserRegisteredEvent) (string, error)
}
