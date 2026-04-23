package ports

import "github.com/sergeyptv/post_service/internal/notification/domain"

type EventPoller interface {
	Poll() (domain.UserRegisteredEvent, error)
	Commit() error
}
