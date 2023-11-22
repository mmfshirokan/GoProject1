package service

import (
	"context"

	"github.com/mmfshirokan/GoProject1/internal/repository"
)

type Password struct {
	rep repository.PwRepositoryInterface
}

func NewPassword(repo repository.PwRepositoryInterface) *Password {
	return &Password{
		rep: repo,
	}
}

func (usr *Password) Store(ctx context.Context, id int, pw string) error {
	return usr.rep.Store(ctx, id, pw)
}

func (usr *Password) Compare(ctx context.Context, id int, pw string) (bool, error) {
	return usr.rep.Compare(ctx, id, pw)
}

func (usr *Password) DeletePassword(ctx context.Context, id int) error {
	return usr.rep.DeletePassword(ctx, id)
}
