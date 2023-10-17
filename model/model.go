package model

type User struct {
	Id   string `param:"id"`
	Name string `query:"name"`
	Male string `query:"male"`
}
