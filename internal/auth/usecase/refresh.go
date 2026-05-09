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

	jti, tokenUser, err := a.tokenSigner.Parse(staleRefreshToken)
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

	_, err = a.sessionRepo.GetToken(ctx, jti)
	if err != nil {
		if errors.Is(err, repository.ErrDbClientClosed) {
			log.Error("Session db client is closed", logger.Error(err))

			return "", "", fmt.Errorf("%s: %w", op, domain.ErrClientNotRespond)
		}

		log.Error("Failed to get refresh token from db", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = a.sessionRepo.DeleteToken(ctx, jti)
	if err != nil {
		log.Error("Failed to delete refresh token from db", logger.Error(err))

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

	_, err = a.sessionRepo.SetToken(ctx, refreshTokenJti, refreshToken, a.config.Redis.TokenTtl)
	if err != nil {
		log.Error("Failed to set token to db", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}
