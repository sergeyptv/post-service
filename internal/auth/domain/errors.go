package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this username or email already exists")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrInvalidTokenType   = errors.New("invalid token type")
	ErrIssIncorrect       = errors.New("iss is incorrect")
	ErrKidNotSet          = errors.New("kid is not set")
	ErrKidIncorrect       = errors.New("kid is incorrect")
	ErrTokenUseNotSet     = errors.New("token use is not set")
	ErrTokenUseIncorrect  = errors.New("token use is incorrect")
	ErrExpFired           = errors.New("exp time is fired")
	ErrClientNotRespond   = errors.New("client is not responding")
)
