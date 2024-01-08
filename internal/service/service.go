package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
)

//go:generate mockery --dir . --all --output ../handlers/mocks/ --with-expecter

type UserInterface interface {
	GetTroughID(ctx context.Context, id int) (usr *model.User, err error)
	Create(ctx context.Context, usr model.User) error
	Update(ctx context.Context, usr model.User) error
	Delete(ctx context.Context, id int) error
}

type User struct {
	rep       repository.RepositoryInterface
	redis     repository.RedisRepositoryInterface[*model.User]
	sourceMap repository.MapRepositoryInterface[*model.User]
}

func NewUser(
	repo repository.RepositoryInterface,
	redis repository.RedisRepositoryInterface[*model.User],
	sourceMap repository.MapRepositoryInterface[*model.User],
) UserInterface {
	return &User{
		rep:       repo,
		redis:     redis,
		sourceMap: sourceMap,
	}
}

func (serv *User) GetTroughID(ctx context.Context, id int) (usr *model.User, err error) {
	key := "user:" + strconv.FormatInt(int64(id), 10)

	usr, err = serv.sourceMap.Get(ctx, key)
	if err != nil {
		usr, err = serv.rep.GetTroughID(ctx, id)
		if err != nil {
			return nil, err
		}

		err = serv.redis.Add(ctx, key, usr)
		if err != nil {
			return usr, fmt.Errorf("redis XAdd error att repository redis: %w", err) // TODO wrap error
		}
	}

	return usr, nil
}

func (serv *User) Create(ctx context.Context, usr model.User) error {
	return serv.rep.Create(ctx, usr)
}

func (serv *User) Update(ctx context.Context, usr model.User) error {
	serv.sourceMap.Remove("user:" + strconv.FormatInt(int64(usr.ID), 10)) // TODO change remove to xadd

	return serv.rep.Update(ctx, usr)
}

func (serv *User) Delete(ctx context.Context, id int) error {
	serv.sourceMap.Remove("user:" + strconv.FormatInt(int64(id), 10))

	return serv.rep.Delete(ctx, id)
}
