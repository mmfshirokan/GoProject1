package service

import (
	"github.com/mmfshirokan/GoProject1/repository"
)

type Password struct {
	rep repository.PwRepositoryInterface
}

func NewPassword(repo repository.PwRepositoryInterface) *Password {
	return &Password{
		rep: repo,
	}
}

func (usr *Password) Store(id uint, pw string) error {
	return usr.rep.Store(id, pw)
}

func (usr *Password) Compare(id uint, pw string) (bool, error) {
	return usr.rep.Compare(id, pw)
}

func (usr *Password) DeletePassword(id uint) error {
	return usr.rep.DeletePassword(id)
}
