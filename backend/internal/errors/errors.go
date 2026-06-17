package errors

import "errors"

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")

	ErrUserNotFound = errors.New("user not found")
	ErrUnauthorized = errors.New("incorrect password")
	ErrTokenWasExpired = errors.New("token was expired")
)
