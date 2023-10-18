package service

import "github.com/mmfshirokan/GoProject1/repository"

type User struct {
	rep repository.RepositoryInterface
}

func NewUser(repo repository.RepositoryInterface) *User {
	return &User{
		rep: repo,
	}
}

func (serv *User) GetTroughID(id int) (string, bool, error) {
	return serv.rep.GetTroughID(id)
}

func (serv *User) Create(id int, name string, male bool) error {
	return serv.rep.Create(id, name, male)
}

func (serv *User) Update(id int, name string, male bool) error {
	return serv.rep.Update(id, name, male)
}

func (serv *User) Delete(id int) error {
	return serv.rep.Delete(id)
}
