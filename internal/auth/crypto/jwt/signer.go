package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	platformJwt "github.com/sergeyptv/post_service/internal/platform/jwt"
	"time"
)

type jwtTokenSigner struct {
	jwtCfg platformJwt.ConfigSigner
}

func NewJwtTokenSigner(jwtCfg platformJwt.ConfigSigner) *jwtTokenSigner {
	return &jwtTokenSigner{
		jwtCfg: jwtCfg,
	}
}

func (j *jwtTokenSigner) NewToken(userUuid, username, userEmail string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, platformJwt.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.jwtCfg.Issuer,
			Subject:   userUuid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.jwtCfg.TokenTTL)),
		},
		Username: username,
		Email:    userEmail,
	},
	)

	return t.SignedString(j.jwtCfg.PrivateKey)
}
