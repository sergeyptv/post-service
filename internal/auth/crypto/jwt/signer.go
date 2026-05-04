package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/sergeyptv/post_service/internal/platform/config"
	platformJwt "github.com/sergeyptv/post_service/internal/platform/jwt"
	"time"
)

type jwtTokenSigner struct {
	app config.App
	cfg Config
}

func NewJwtTokenSigner(cfg Config) *jwtTokenSigner {
	return &jwtTokenSigner{cfg: cfg}
}

func (j *jwtTokenSigner) NewToken(userUuid, username, userEmail string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, platformJwt.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.app.Name,
			Subject:   userUuid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.TokenTTL)),
		},
		Username: username,
		Email:    userEmail,
	},
	)

	return t.SignedString(j.cfg.PrivateKey)
}
