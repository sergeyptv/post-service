package ports

import (
	"github.com/sergeyptv/post_service/auth/internal/domain"
)

type TokenSigner interface {
	NewToken(userUuid, username, userEmail, tokenType string) (jti string, signedToken string, err error)
	Parse(jwtToken, tokenType string) (jti string, user domain.User, err error)
}
