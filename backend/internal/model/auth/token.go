package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
	UserId       uuid.UUID
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type Claims struct {
	UserId uuid.UUID
	jwt.RegisteredClaims
}

func NewAccessClaims(userId uuid.UUID) Claims {
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return claims
}

func NewRefreshClaims(userId uuid.UUID) Claims {
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return claims
}
