package service

import "github.com/mmfshirokan/GoProject1/repository"

type User struct {
	rep *repository.Repository
}

func NewService(repo *repository.Repository) *User {
	return &User{
		rep: repo,
	}
}

func (serv *User) GetTroughID(id string) (string, string, error) {
	return serv.rep.GetUserTroughID(id)
}

func (serv *User) Create(id string, name string, male string) error {
	return serv.rep.SaveUser(id, name, male)
}

func (serv *User) Update(id string, name string, male string) error {
	return serv.rep.UpdateUser(id, name, male)
}

func (serv *User) Delete(id string) error {
	return serv.rep.DeleteUser(id)
}
