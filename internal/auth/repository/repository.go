package repository

import "errors"

var (
	ErrUserExists   = errors.New("user with this username, email or phone exists")
	ErrUserNotFound = errors.New("user is not found")
)
