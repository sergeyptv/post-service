package repository

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with this username or email already exists")
	ErrUserNotFound      = errors.New("user is not found")
	ErrNoRowsAffected    = errors.New("no rows affected")
	ErrTokenNotFound     = errors.New("token is not found")
)
