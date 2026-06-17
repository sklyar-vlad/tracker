package errors

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailNotFound = errors.New("email not found")
)
