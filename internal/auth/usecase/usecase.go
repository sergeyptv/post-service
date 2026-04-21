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
	idemRepo    ports.IdempotencyRepository
	tokenSigner ports.TokenSigner
	txWrapper   ports.TransactionWrapper
}

func NewAuthService(log *slog.Logger, userRepo ports.UserRepository, outboxRepo ports.OutboxRepository, tokenRepo ports.TokenRepository, idemRepo ports.IdempotencyRepository, tokenSigner ports.TokenSigner, txWrapper ports.TransactionWrapper) *auth {
	return &auth{
		log:         log,
		userRepo:    userRepo,
		outboxRepo:  outboxRepo,
		tokenRepo:   tokenRepo,
		idemRepo:    idemRepo,
		tokenSigner: tokenSigner,
		txWrapper:   txWrapper,
	}
}
