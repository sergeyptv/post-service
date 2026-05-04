package jwt

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	platformJwt "github.com/sergeyptv/post_service/internal/platform/jwt"
	"github.com/sergeyptv/post_service/internal/post/domain"
	"github.com/sergeyptv/post_service/internal/post/ports"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	ErrIssIncorrect = errors.New("iss is incorrect")
	ErrExpFired     = errors.New("exp time is fired")
	ErrTokenInvalid = errors.New("token is invalid")
	ErrGetPublicKey = errors.New("error getting public key")
)

type jwtTokenParser struct {
	jwtCache   JwtCache
	authClient ports.AuthClient
	sf         singleflight.Group
}

func NewJwtTokenParser(jwtCache JwtCache, authClient ports.AuthClient) *jwtTokenParser {
	return &jwtTokenParser{
		jwtCache:   jwtCache,
		authClient: authClient,
	}
}

func (j *jwtTokenParser) publicKey(ctx context.Context) (*rsa.PublicKey, error) {
	const op = "auth.jwt.publicKey"

	key, err := j.jwtCache.Get()
	if err == nil {
		return key, nil
	}

	v, err, _ := j.sf.Do("publicKey", func() (any, error) {
		key, err := j.jwtCache.Get()
		if err == nil {
			return key, nil
		}

		return j.refreshKey(ctx)
	})
	if err != nil {
		key, cacheErr := j.jwtCache.Get()
		if cacheErr == nil {
			return key, nil
		}

		if errors.Is(cacheErr, ErrPublicKeyTtlExpired) {
			return key, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	key, ok := v.(*rsa.PublicKey)
	if !ok || key == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrGetPublicKey)
	}

	return key, nil
}

func (j *jwtTokenParser) refreshKey(ctx context.Context) (*rsa.PublicKey, error) {
	const op = "auth.jwt.refreshKey"

	rsaPublicKey, err := j.authClient.GetPublicKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	j.jwtCache.Set(rsaPublicKey)

	return rsaPublicKey, nil
}

func (j *jwtTokenParser) validate(claims platformJwt.Claims) error {
	const op = "auth.jwt.validate"

	iss, err := claims.GetIssuer()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if iss != "auth" {
		return fmt.Errorf("%s: %w", op, ErrIssIncorrect)
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if exp.Before(time.Now()) {
		return fmt.Errorf("%s: %w", op, ErrExpFired)
	}

	return nil
}

func (j *jwtTokenParser) Parse(ctx context.Context, jwtToken string) (domain.User, error) {
	const op = "auth.jwt.validate"

	var claims platformJwt.Claims

	token, err := jwt.ParseWithClaims(jwtToken, &claims, func(*jwt.Token) (any, error) {
		return j.publicKey(ctx)
	}, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	if !token.Valid {
		return domain.User{}, fmt.Errorf("%s: %w", op, ErrTokenInvalid)
	}

	err = j.validate(claims)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return domain.User{
		Uuid:     sub,
		Username: claims.Username,
		Email:    claims.Email,
	}, nil
}
