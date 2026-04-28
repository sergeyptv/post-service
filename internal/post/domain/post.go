package domain

import "time"

type Post struct {
	Uuid        string    `json:"uuid"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Media       []byte    `json:"media"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
