package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)
