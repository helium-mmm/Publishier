package domain

import "errors"

var (
	ErrPostNotFound       = errors.New("post not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidStatus      = errors.New("invalid post status")
	ErrAccountNotFound    = errors.New("telegram channel not connected")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrPublicationFailed  = errors.New("publication failed")
)
