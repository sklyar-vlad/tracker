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

type AuthResponse struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

func ToUserResponse(p model.User) UserResponse {
	resp := UserResponse{
		User_id:  p.User_id,
		Email:    p.Email,
		Password: p.Password,
	}

	return resp
}

func ToAuthResponse(access_token, refresh_token string) AuthResponse {
	resp := AuthResponse{
		Access_token:  access_token,
		Refresh_token: refresh_token,
	}

	return resp
}
