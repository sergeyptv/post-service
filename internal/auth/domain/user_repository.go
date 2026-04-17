package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUser) (string, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
