package service

import (
	"context"

	"github.com/mmfshirokan/GoProject1/internal/model"
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

func (usr *Password) Store(ctx context.Context, user model.User) error {
	return usr.rep.Store(ctx, user)
}

func (usr *Password) Compare(ctx context.Context, user model.User) (bool, error) {
	return usr.rep.Compare(ctx, user)
}

func (usr *Password) DeletePassword(ctx context.Context, id int) error {
	return usr.rep.DeletePassword(ctx, id)
}
