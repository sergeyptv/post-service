package usecase

import (
	"context"
	"fmt"
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

	userUuid, err := a.UserRepo.CreateUser(ctx, user)
	if err != nil {
		log.Error("Failed to create user", logger.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		for i := 0; i < 5; i++ {
			err = a.Publisher.Publish(ctx, "user.registered", domain.UserRegisteredEvent{
				Version:      "v1",
				UserUuid:     userUuid,
				UserEmail:    user.Email,
				RegisteredAt: time.Now().Unix(),
			})
			if err != nil {
				log.Error("Failed to publish new user", logger.Error(err))
			} else {
				return
			}
		}

		err = a.Publisher.Publish(ctx, "user.registered.dlq", domain.UserRegisteredEvent{
			Version:      "v1",
			UserUuid:     userUuid,
			UserEmail:    user.Email,
			RegisteredAt: time.Now().Unix(),
		})
		if err != nil {
			log.Error("Failed to publish new user", logger.Error(err))
		}
	}()

	return userUuid, nil
}
