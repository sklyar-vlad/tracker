package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	TokenHash string
	User_id  uuid.UUID
	CreatedAt time.Time
	ExpiresAt time.Time
}
