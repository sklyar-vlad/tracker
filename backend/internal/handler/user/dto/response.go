package dto

import (
	"github.com/google/uuid"
	"github.com/sklyar-vlad/selfDev/internal/model"
)

type UserResponse struct {
	User_id  uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func ToUserResponse(p model.User) UserResponse {
	resp := UserResponse{
		User_id:  p.User_id,
		Email:    p.Email,
		Password: p.Password,
	}

	return resp
}
