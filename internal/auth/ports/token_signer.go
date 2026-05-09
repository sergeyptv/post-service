package ports

import "github.com/sergeyptv/post_service/internal/auth/domain"

type TokenSigner interface {
	NewToken(userUuid, username, userEmail, tokenType string) (jti string, signedToken string, err error)
	Parse(jwtToken string) (jti string, user domain.User, err error)
}
