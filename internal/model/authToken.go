package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Male bool   `json:"male"`
	jwt.RegisteredClaims
}
