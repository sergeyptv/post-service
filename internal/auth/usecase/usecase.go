package usecase

import (
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"log/slog"
)

type auth struct {
	log         *slog.Logger
	UserRepo    domain.UserRepository
	TokenRepo   domain.TokenRepository
	Publisher   domain.EventPublisher
	IdemRepo    domain.IdempotencyRepository
	TokenSigner domain.TokenSigner
}

func NewAuthService(log *slog.Logger, userRepo domain.UserRepository, tokenRepo domain.TokenRepository, publisher domain.EventPublisher, idemRepo domain.IdempotencyRepository, tokenSigner domain.TokenSigner) *auth {
	return &auth{
		log:         log,
		UserRepo:    userRepo,
		TokenRepo:   tokenRepo,
		Publisher:   publisher,
		IdemRepo:    idemRepo,
		TokenSigner: tokenSigner,
	}
}
