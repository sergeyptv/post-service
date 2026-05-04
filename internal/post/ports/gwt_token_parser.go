package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/post/domain"
)

type JwtTokenParser interface {
	Parse(ctx context.Context, jwtToken string) (domain.User, error)
}
