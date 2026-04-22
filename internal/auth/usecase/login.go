package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (a *auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "usecase.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("log in user")

	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Error("User not found", logger.Error(err))

			return "", fmt.Errorf("%s: %w", op, domain.ErrInvalidCredentials)
		}

		log.Error("Failed to get user", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		log.Error("Invalid credentials", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := a.tokenSigner.NewToken(user.Uuid, user.Email)
	if err != nil {
		log.Error("Failed to create token", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	return token, nil
}
