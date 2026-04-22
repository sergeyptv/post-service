package usecase

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/platform/logger"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

func (a *auth) Register(ctx context.Context, user domain.CreateUser) (string, error) {
	const op = "usecase.Register"

	log := a.log.With(slog.String("op", op), slog.String("email", user.Email))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.PassHash), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to generate password hash", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	user.PassHash = string(passHash)

	var userUuid string

	err = a.txWrapper.Wrap(ctx, func(ctx context.Context, tx pgx.Tx) error {
		userUuid, terr := a.userRepo.CreateUser(ctx, tx, user)
		if terr != nil {
			return terr
		}

		_, terr = a.outboxRepo.CreateEvent(ctx,
			tx,
			domain.UserRegisteredEvent{
				Version:      "1.0",
				UserUuid:     userUuid,
				Username:     user.Username,
				UserEmail:    user.Email,
				RegisteredAt: time.Now().UTC(),
			})
		if terr != nil {
			return terr
		}

		return nil
	})
	if err != nil {
		log.Error("Failed to add userinfo to db", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userUuid, nil
}
