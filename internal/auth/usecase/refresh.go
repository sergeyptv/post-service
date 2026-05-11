package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/auth/crypto/jwt"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"log/slog"
)

func (a *auth) Refresh(ctx context.Context, staleRefreshToken string) (accessToken string, refreshToken string, err error) {
	const op = "usecase.Refresh"

	log := a.log.With(slog.String("op", op))

	staleJti, tokenUser, err := a.tokenSigner.Parse(staleRefreshToken)
	if err != nil {
		log.Error("Failed to parse refresh token", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.userRepo.GetUserByEmail(ctx, tokenUser.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Error("User not found", logger.Error(err))

			return "", "", fmt.Errorf("%s: %w", op, domain.ErrTokenInvalid)
		}

		log.Error("Failed to get user", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	_, accessToken, err = a.tokenSigner.NewToken(user.Uuid, user.Username, user.Email, jwt.TypeAccess)
	if err != nil {
		log.Error("Failed to create access token", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refreshTokenJti, refreshToken, err := a.tokenSigner.NewToken(user.Uuid, user.Username, user.Email, jwt.TypeRefresh)
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
