package service

import (
	"context"
	"fmt"

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

func (usr *Password) Store(ctx context.Context, id int, pw string) error {
	return fmt.Errorf("rep.Store: %w", usr.rep.Store(ctx, id, pw))
}

func (usr *Password) Compare(ctx context.Context, id int, pw string) (bool, error) {
	bl, err := usr.rep.Compare(ctx, id, pw)
	return bl, fmt.Errorf("rep.Compare: %w", err)
}

func (usr *Password) DeletePassword(ctx context.Context, id int) error {
	return fmt.Errorf("rep.DeletePassword: %w", usr.rep.DeletePassword(ctx, id))
}
