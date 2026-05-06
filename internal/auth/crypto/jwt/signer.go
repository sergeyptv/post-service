package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/sergeyptv/post_service/internal/platform/config"
	platformJwt "github.com/sergeyptv/post_service/internal/platform/jwt"
	"time"
)

type jwtTokenSigner struct {
	appCfg config.App
	jwtCfg Config
}

func NewJwtTokenSigner(appCfg config.App, jwtCfg Config) *jwtTokenSigner {
	return &jwtTokenSigner{
		appCfg: appCfg,
		jwtCfg: jwtCfg,
	}
}

func (j *jwtTokenSigner) NewToken(userUuid, username, userEmail string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, platformJwt.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.appCfg.Name,
			Subject:   userUuid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.jwtCfg.TokenTTL)),
		},
		Username: username,
		Email:    userEmail,
	},
	)

	return t.SignedString(j.jwtCfg.PrivateKey)
}
