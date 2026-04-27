package repository

import "errors"

var (
	ErrEventStatusProcessing = errors.New("event status processing")
	ErrEventAlreadySuccess   = errors.New("event status already success")
)
