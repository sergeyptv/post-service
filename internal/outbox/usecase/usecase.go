package usecase

import (
	"github.com/sergeyptv/post_service/internal/outbox/ports"
	"log/slog"
)

type outbox struct {
	log              *slog.Logger
	outboxRepository ports.OutboxRepository
	publisher        ports.EventPublisher
	txWrapper        ports.TransactionWrapper
}

func NewOutboxService(log *slog.Logger, outboxRepository ports.OutboxRepository, publisher ports.EventPublisher, txWrapper ports.TransactionWrapper) *outbox {
	return &outbox{
		log:              log,
		outboxRepository: outboxRepository,
		publisher:        publisher,
		txWrapper:        txWrapper,
	}
}
