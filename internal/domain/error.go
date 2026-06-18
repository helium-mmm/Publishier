package domain

import "errors"

var (
	ErrPostNotFound = errors.New("post not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidStatus = errors.New("invalid post status")
	ErrAccountNotFound = errors.New("social account not found")
)