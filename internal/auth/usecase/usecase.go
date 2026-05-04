package usecase

import (
	"github.com/sergeyptv/post_service/internal/auth/ports"
	"log/slog"
)

type auth struct {
	log         *slog.Logger
	userRepo    ports.UserRepository
	outboxRepo  ports.OutboxRepository
	tokenRepo   ports.TokenRepository
	tokenSigner ports.TokenSigner
	txWrapper   ports.TransactionWrapper
}

func NewAuthService(log *slog.Logger, userRepo ports.UserRepository, outboxRepo ports.OutboxRepository, tokenRepo ports.TokenRepository, tokenSigner ports.TokenSigner, txWrapper ports.TransactionWrapper) *auth {
	return &auth{
		log:         log,
		userRepo:    userRepo,
		outboxRepo:  outboxRepo,
		tokenRepo:   tokenRepo,
		tokenSigner: tokenSigner,
		txWrapper:   txWrapper,
	}
}
