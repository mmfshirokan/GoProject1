package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID     int       `json:"UserID"`
	ID         uuid.UUID `json:"ID"`
	Hash       string    `json:"Hash"`
	Expiration time.Time `json:"Expiration"`
}
