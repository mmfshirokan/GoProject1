package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID     int
	ID         uuid.UUID
	Hash       string
	Expiration time.Time
}
