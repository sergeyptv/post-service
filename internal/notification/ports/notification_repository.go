package ports

import (
	"context"
	"github.com/sergeyptv/post_service/internal/notification/domain"
)

type NotificationRepository interface {
	GetStaleEvent(ctx context.Context) (domain.UserRegisteredEvent, error)
	TryProcess(ctx context.Context, event domain.UserRegisteredEvent) error
	MarkSuccess(ctx context.Context, eventUuid string) error
}
