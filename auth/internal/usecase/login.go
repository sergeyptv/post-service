package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/auth/internal/crypto/jwt"
	"github.com/sergeyptv/post_service/auth/internal/domain"
	"github.com/sergeyptv/post_service/auth/internal/repository"
	"github.com/sergeyptv/post_service/platform/logger"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (a *auth) Login(ctx context.Context, email, password string) (string, string, error) {
	const op = "usecase.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Warn("User not found", logger.Error(err))

			return "", "", fmt.Errorf("%s: %w", op, domain.ErrInvalidCredentials)
		}

		log.Error("Failed to get user", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Warn("Invalid credentials", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, domain.ErrInvalidCredentials)
	}

	_, accessToken, err := a.tokenSigner.NewToken(user.Uuid, user.Username, user.Email, jwt.TypeAccess)
	if err != nil {
		log.Error("Failed to create access token", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refreshTokenJti, refreshToken, err := a.tokenSigner.NewToken(user.Uuid, user.Username, user.Email, jwt.TypeRefresh)
	if err != nil {
		log.Error("Failed to create refresh token", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = a.sessionRepo.SetToken(ctx, fmt.Sprintf("sess:%s", refreshTokenJti), refreshToken, a.config.Redis.TokenTtl)
	if err != nil {
		log.Error("Failed to set token to db", logger.Error(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}
