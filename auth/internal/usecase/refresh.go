package usecase

import (
	"context"
	"fmt"
	"github.com/sergeyptv/post_service/auth/internal/crypto/jwt"
	"github.com/sergeyptv/post_service/auth/internal/domain"
	"github.com/sergeyptv/post_service/platform/logger"
	"log/slog"
)

func (a *auth) Refresh(ctx context.Context, staleRefreshToken string) (accessToken string, refreshToken string, err error) {
	const op = "usecase.Refresh"

	log := a.log.With(slog.String("op", op))

	staleJti, tokenUser, err := a.tokenSigner.Parse(staleRefreshToken, jwt.TypeRefresh)
	if err != nil {
		log.Error("Failed to parse refresh token", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	_, accessToken, err = a.tokenSigner.NewToken(tokenUser.Uuid, tokenUser.Username, tokenUser.Email, jwt.TypeAccess)
	if err != nil {
		log.Error("Failed to create access token", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refreshTokenJti, refreshToken, err := a.tokenSigner.NewToken(tokenUser.Uuid, tokenUser.Username, tokenUser.Email, jwt.TypeRefresh)
	if err != nil {
		log.Error("Failed to create refresh token", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	ok, err := a.sessionRepo.RotateToken(ctx, fmt.Sprintf("sess:%s", staleJti), fmt.Sprintf("sess:%s", refreshTokenJti), refreshToken, a.config.Redis.TokenTtl)
	if err != nil {
		log.Error("Failed to rotate token in db", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if !ok {
		log.Warn("staleJti is invalid", logger.Error(domain.ErrTokenInvalid))

		return "", "", fmt.Errorf("%s: %w", op, domain.ErrTokenInvalid)
	}

	return accessToken, refreshToken, nil
}
