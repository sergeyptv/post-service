package domain

type UserRegisteredEvent struct {
	Version      string `json:"version"`
	UserUuid     string `json:"user_uuid"`
	UserEmail    string `json:"user_email"`
	RegisteredAt int64  `json:"registered_at"`
}
