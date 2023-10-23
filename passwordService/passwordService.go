package passwordService

import (
	"github.com/mmfshirokan/GoProject1/passwordRepository"
)

type Password struct {
	rep passwordRepository.RepositoryInterface
}

func NewPassword(repo passwordRepository.RepositoryInterface) *Password {
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
