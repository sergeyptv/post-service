package domain

import "time"

type UserRegisteredEvent struct {
	Version      string    `json:"version"`
	UserUuid     string    `json:"user_uuid"`
	Username     string    `json:"username"`
	UserEmail    string    `json:"user_email"`
	RegisteredAt time.Time `json:"registered_at"`
}
