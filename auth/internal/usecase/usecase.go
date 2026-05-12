package usecase

import (
	"github.com/sergeyptv/post_service/auth/internal/config"
	"github.com/sergeyptv/post_service/auth/internal/ports"
	"log/slog"
)

type auth struct {
	log         *slog.Logger
	config      *config.Config
	userRepo    ports.UserRepository
	outboxRepo  ports.OutboxRepository
	sessionRepo ports.SessionRepository
	tokenSigner ports.TokenSigner
	txWrapper   ports.TransactionWrapper
}

func NewAuthUsecase(log *slog.Logger, config *config.Config, userRepo ports.UserRepository, outboxRepo ports.OutboxRepository, sessionRepo ports.SessionRepository, tokenSigner ports.TokenSigner, txWrapper ports.TransactionWrapper) *auth {
	return &auth{
		log:         log,
		config:      config,
		userRepo:    userRepo,
		outboxRepo:  outboxRepo,
		sessionRepo: sessionRepo,
		tokenSigner: tokenSigner,
		txWrapper:   txWrapper,
	}
}
