package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	platformJwt "github.com/sergeyptv/post_service/internal/platform/jwt"
	"time"
)

const (
	TypeAccess  = "access"
	TypeRefresh = "refresh"
)

type jwtTokenSigner struct {
	config platformJwt.ConfigSigner
}

func NewJwtTokenSigner(config platformJwt.ConfigSigner) *jwtTokenSigner {
	return &jwtTokenSigner{
		config: config,
	}
}

func (j *jwtTokenSigner) NewToken(userUuid, username, userEmail, tokenType string) (jti string, signedToken string, err error) {
	var ttl time.Duration

	switch tokenType {
	case TypeAccess:
		ttl = j.config.AccessTokenTtl

	case TypeRefresh:
		ttl = j.config.RefreshTokenTtl

	default:
		return "", "", domain.ErrInvalidTokenType
	}

	jti = uuid.New().String()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, platformJwt.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Subject:   userUuid,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			ID:        jti,
		},
		Username: username,
		Email:    userEmail,
	})

	t.Header["kid"] = j.config.Kid
	t.Header["token_use"] = tokenType

	signedToken, err = t.SignedString(j.config.PrivateKey)
	if err != nil {
		return "", "", err
	}

	return jti, signedToken, nil
}

func (j *jwtTokenSigner) Parse(jwtToken string) (jti string, user domain.User, err error) {
	var claims platformJwt.Claims

	token, err := jwt.ParseWithClaims(jwtToken, &claims, func(*jwt.Token) (any, error) {
		return j.config.PublicKey, nil
	}, jwt.WithValidMethods([]string{j.config.Algorithm}))
	if err != nil {
		return "", domain.User{}, err
	}

	if !token.Valid {
		return "", domain.User{}, domain.ErrTokenInvalid
	}

	err = j.validateHeader(token.Header)
	if err != nil {
		return "", domain.User{}, err
	}

	err = j.validate(claims)
	if err != nil {
		return "", domain.User{}, err
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return "", domain.User{}, err
	}

	return claims.ID, domain.User{
		Uuid:     sub,
		Username: claims.Username,
		Email:    claims.Email,
	}, nil
}

func (j *jwtTokenSigner) validateHeader(header map[string]any) error {
	kid, ok := header["kid"]
	if !ok || kid == "" {
		return domain.ErrKidNotSet
	}
	if kid != j.config.Kid {
		return domain.ErrKidIncorrect
	}

	tokenUse, ok := header["token_use"]
	if !ok || tokenUse == "" {
		return domain.ErrTokenUseNotSet
	}
	if tokenUse != TypeRefresh {
		return domain.ErrTokenUseIncorrect
	}

	return nil
}

func (j *jwtTokenSigner) validate(claims platformJwt.Claims) error {
	iss, err := claims.GetIssuer()
	if err != nil {
		return err
	}
	if iss != j.config.Issuer {
		return domain.ErrIssIncorrect
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return err
	}
	if exp.Before(time.Now()) {
		return domain.ErrExpFired
	}

	return nil
}
