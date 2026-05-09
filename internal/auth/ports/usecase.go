package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/auth/domain"
)

type Usecase interface {
	Register(ctx context.Context, user domain.User, password string) (userUuid string, err error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	Refresh(ctx context.Context, staleRefreshToken string) (accessToken string, refreshToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
}
