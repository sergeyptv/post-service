package domain

import "time"

type UserRegisteredEvent struct {
	Uuid         string
	Version      string    `json:"version"`
	UserUuid     string    `json:"user_uuid"`
	UserEmail    string    `json:"user_email"`
	RegisteredAt time.Time `json:"registered_at"`
}
