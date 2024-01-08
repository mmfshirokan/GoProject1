package service

import (
	"context"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
)

type PasswordInterface interface {
	Store(ctx context.Context, usr model.User) error
	Compare(ctx context.Context, usr model.User) (bool, error)
	DeletePassword(ctx context.Context, id int) error
}

type Password struct {
	rep repository.PwRepositoryInterface
}

func NewPassword(repo repository.PwRepositoryInterface) PasswordInterface {
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
