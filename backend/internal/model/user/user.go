package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/sklyar-vlad/selfDev/internal/errors"
)

type User struct {
	UserId        uuid.UUID
	Role          string
	Username      string
	Email         string
	EmailVerified bool
	Password      string
}

func NewUser(username, email, password string) (User, error) {
	if !strings.ContainsRune(email, '@') {
		return User{}, fmt.Errorf("invalid email: %w", errors.ErrInvalidEmail)
	}

	if len(password) < 6 {
		return User{}, fmt.Errorf("invalid password: %w", errors.ErrInvalidPassword)
	}

	return User{uuid.New(), "user", username, email, false, password}, nil
}
