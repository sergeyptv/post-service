package usecase

import (
	"context"
	"fmt"
	"github.com/sergeyptv/post_service/auth/internal/crypto/jwt"
	"github.com/sergeyptv/post_service/platform/logger"
	"log/slog"
)

func (a *auth) Logout(ctx context.Context, refreshToken string) error {
	const op = "usecase.Logout"

	log := a.log.With(slog.String("op", op))

	jti, _, err := a.tokenSigner.Parse(refreshToken, jwt.TypeRefresh)
	if err != nil {
		log.Error("Failed to parse refresh token", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.sessionRepo.DeleteToken(ctx, jti)
	if err != nil {
		log.Error("Failed to delete token from db", logger.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
