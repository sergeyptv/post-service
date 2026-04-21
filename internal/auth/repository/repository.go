package repository

import "errors"

var (
	ErrUserExists   = errors.New("user with this username or email already exists")
	ErrUserNotFound = errors.New("user is not found")
)
