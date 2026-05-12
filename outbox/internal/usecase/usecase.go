package usecase

import (
	"github.com/sergeyptv/post_service/outbox/internal/config"
	"github.com/sergeyptv/post_service/outbox/internal/ports"
	"log/slog"
)

type outbox struct {
	log              *slog.Logger
	cfg              *config.Config
	outboxRepository ports.OutboxRepository
	publisher        ports.EventPublisher
	txWrapper        ports.TransactionWrapper
}

func NewOutboxService(log *slog.Logger, cfg *config.Config, outboxRepository ports.OutboxRepository, publisher ports.EventPublisher, txWrapper ports.TransactionWrapper) *outbox {
	return &outbox{
		log:              log,
		cfg:              cfg,
		outboxRepository: outboxRepository,
		publisher:        publisher,
		txWrapper:        txWrapper,
	}
}
