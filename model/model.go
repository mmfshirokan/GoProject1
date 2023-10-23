package model

type User struct {
	Id   uint   `param:"id" query:"id"`
	Name string `query:"name"`
	Male bool   `query:"male"`
}
