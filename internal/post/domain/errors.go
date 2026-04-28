package domain

import "errors"

var (
	ErrBadGateway   = errors.New("bad gateway")
	ErrPostNotExist = errors.New("post does not exist")
)
