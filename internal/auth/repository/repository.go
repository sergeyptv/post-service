package repository

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with this username or email already exists")
	ErrUserNotFound      = errors.New("user is not found")
	ErrDbClientClosed    = errors.New("db client closed")
)
