package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"log/slog"
)

func (a *auth) Logout(ctx context.Context, refreshToken string) error {
	const op = "usecase.Logout"

	log := a.log.With(slog.String("op", op))

	jti, tokenUser, err := a.tokenSigner.Parse(refreshToken)
	if err != nil {
		log.Error("Failed to parse refresh token", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = a.userRepo.GetUserByEmail(ctx, tokenUser.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Error("User not found", logger.Error(err))

			return fmt.Errorf("%s: %w", op, domain.ErrTokenInvalid)
		}

		log.Error("Failed to get user", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.sessionRepo.Delete(ctx, jti)
	if err != nil {
		log.Error("Failed to delete token from db", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
