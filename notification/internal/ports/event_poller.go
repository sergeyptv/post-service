package ports

import (
	"github.com/sergeyptv/post_service/notification/internal/domain"
)

type EventPoller interface {
	Poll() (domain.UserRegisteredEvent, error)
	Commit() error
}
