package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserRequest struct {
	ID   int    `json:"id" validate:"max=1000000000"`
	Name string `json:"name" validate:"max=40"`
	Male bool   `json:"male" validate:"boolean"`
	jwt.RegisteredClaims
}
