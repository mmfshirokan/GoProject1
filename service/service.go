package service

import (
	"context"

	"github.com/mmfshirokan/GoProject1/repository"
)

type User struct {
	rep repository.RepositoryInterface
}

func NewUser(repo repository.RepositoryInterface) *User {
	return &User{
		rep: repo,
	}
}

func (serv *User) GetTroughID(ctx context.Context, id int) (string, bool, error) {
	return serv.rep.GetTroughID(ctx, id)
}

func (serv *User) Create(ctx context.Context, id int, name string, male bool) error {
	return serv.rep.Create(ctx, id, name, male)
}

func (serv *User) Update(ctx context.Context, id int, name string, male bool) error {
	return serv.rep.Update(ctx, id, name, male)
}

func (serv *User) Delete(ctx context.Context, id int) error {
	return serv.rep.Delete(ctx, id)
}
