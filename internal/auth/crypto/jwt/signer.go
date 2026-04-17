package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtTokenSigner struct {
	cfg Config
}

func NewJwtTokenSigner(cfg Config) *jwtTokenSigner {
	return &jwtTokenSigner{cfg: cfg}
}

func (j *jwtTokenSigner) NewToken(userUuid, userEmail string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "auth",
		"sub":   userUuid,
		"email": userEmail,
		"exp":   time.Now().Add(j.cfg.TokenTTL).Unix(),
	})

	return t.SignedString(j.cfg.TokenKey)
}
