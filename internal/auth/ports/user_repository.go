package ports

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx pgx.Tx, user domain.CreateUser) (string, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}
