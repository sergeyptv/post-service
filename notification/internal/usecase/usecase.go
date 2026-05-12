package usecase

import (
	"github.com/sergeyptv/post_service/notification/internal/config"
	"github.com/sergeyptv/post_service/notification/internal/ports"
	"log/slog"
)

type notification struct {
	log                    *slog.Logger
	cfg                    *config.Config
	notificationRepository ports.NotificationRepository
	poller                 ports.EventPoller
}

func NewNotificationService(log *slog.Logger, cfg *config.Config, notificationRepository ports.NotificationRepository, poller ports.EventPoller) *notification {
	return &notification{
		log:                    log,
		cfg:                    cfg,
		notificationRepository: notificationRepository,
		poller:                 poller,
	}
}
