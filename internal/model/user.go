package model

type User struct {
	ID       int    `json:"id" validate:"max=1000000000"`
	Name     string `json:"name" validate:"max=40"`
	Male     bool   `json:"male"`
	Password string `json:"password" validate:"max=40"`
}
