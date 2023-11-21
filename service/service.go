package service

import (
	"context"
	"fmt"

	"github.com/mmfshirokan/GoProject1/repository"
)

type User struct {
	rep repository.Interface
}

func NewUser(repo repository.Interface) *User {
	return &User{
		rep: repo,
	}
}

func (serv *User) GetTroughID(ctx context.Context, id int) (string, bool, error) {
	str, bl, err := serv.rep.GetTroughID(ctx, id)

	return str, bl, fmt.Errorf("rep.GetTroughID: %w", err)
}

func (serv *User) Create(ctx context.Context, id int, name string, male bool) error {
	return fmt.Errorf("rep.Create: %w", serv.rep.Create(ctx, id, name, male))
}

func (serv *User) Update(ctx context.Context, id int, name string, male bool) error {
	return fmt.Errorf("rep.Update: %w", serv.rep.Update(ctx, id, name, male))
}

func (serv *User) Delete(ctx context.Context, id int) error {
	return fmt.Errorf("rep.Delete: %w", serv.rep.Delete(ctx, id))
}
