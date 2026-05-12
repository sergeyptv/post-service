package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/auth/internal/domain"
	"github.com/sergeyptv/post_service/auth/internal/repository"
	"github.com/sergeyptv/post_service/platform/logger"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (a *auth) Register(ctx context.Context, user domain.User, password string) (userUuid string, err error) {
	const op = "usecase.Register"

	log := a.log.With(slog.String("op", op), slog.String("email", user.Email))

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to generate password hash", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	user.PasswordHash = string(passwordHash)

	err = a.txWrapper.Wrap(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var terr error

		userUuid, terr = a.userRepo.CreateUser(ctx, tx, user)
		if terr != nil {
			if errors.Is(terr, repository.ErrUserAlreadyExists) {
				return domain.ErrUserAlreadyExists
			}

			return terr
		}

		_, terr = a.outboxRepo.CreateEvent(ctx,
			tx,
			domain.UserRegisteredEvent{
				Version:   "1.0",
				UserUuid:  userUuid,
				Username:  user.Username,
				UserEmail: user.Email,
			})
		if terr != nil {
			return terr
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			log.Warn("Failed to add user info to db", logger.Error(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}
		log.Error("Failed to add user info to db", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userUuid, nil
}
