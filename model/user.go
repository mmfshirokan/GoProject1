package model

type User struct {
	Id       int    `json:"id" param:"id" query:"id"`
	Name     string `json:"name" query:"name"`
	Male     bool   `json:"male" query:"male"`
	Password string `json:"password" query:"password"`
}
