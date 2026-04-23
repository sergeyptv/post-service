package ports

import "context"

type NotificationRepository interface {
	TryProcess(ctx context.Context, eventUuid string) error
	MarkSuccess(ctx context.Context, eventUuid string) error
}
