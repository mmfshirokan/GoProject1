package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID     int       `json:"userId"`
	ID         uuid.UUID `json:"id"`
	Hash       string    `json:"hash"`
	Expiration time.Time `json:"expiration"`
}
