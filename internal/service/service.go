package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
)

type User struct {
	rep       repository.Interface
	redis     *repository.RedisRepository[*model.User]
	sourceMap *repository.MapRepository[*model.User]
}

func NewUser(
	repo repository.Interface,
	redis *repository.RedisRepository[*model.User],
	sourceMap *repository.MapRepository[*model.User],
) *User {
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

func (serv *User) Create(ctx context.Context, id int, name string, male bool) error {
	return serv.rep.Create(ctx, id, name, male)
}

func (serv *User) Update(ctx context.Context, id int, name string, male bool) error {
	serv.sourceMap.Remove("user:" + strconv.FormatInt(int64(id), 10)) // TODO change remove to xadd

	return serv.rep.Update(ctx, id, name, male)
}

func (serv *User) Delete(ctx context.Context, id int) error {
	serv.sourceMap.Remove("user:" + strconv.FormatInt(int64(id), 10))

	return serv.rep.Delete(ctx, id)
}
