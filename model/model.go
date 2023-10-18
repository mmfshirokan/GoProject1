package model

type User struct {
	Id   int    `param:"id"`
	Name string `query:"name"`
	Male bool   `query:"male"`
}
