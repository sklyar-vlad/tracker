package model

import (
	"github.com/google/uuid"
)

type User struct {
	UserId   uuid.UUID
	Role     string
	Username string
	Email    string
	Password string
}

func NewUser(username, email, password string) (User, error) {
	return User{uuid.New(), "user", username, email, password}, nil
}
