package model

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Male     bool   `json:"male"`
	Password string `json:"password"`
}
