package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID     int       `json:"userId" validate:"max=1000000000"`
	ID         uuid.UUID `json:"id" validate:"uuid"`
	Hash       string    `json:"hash"` //validate:"sha256"
	Expiration time.Time `json:"expiration"`
}
