package ports

import (
	"context"
	"github.com/sergeyptv/post_service/post/internal/domain"
)

type JwtTokenParser interface {
	Parse(ctx context.Context, jwtToken, tokenType string) (domain.User, error)
}
