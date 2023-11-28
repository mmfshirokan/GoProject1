package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/redis/go-redis/v9"
)

type repositoryRedis[generic model.User | []*model.RefreshToken] struct {
	client *redis.Client
	model  *generic
}

func NewUserRedisRepository(conf config.Config) *repositoryRedis[model.User] {
	opts, err := redis.ParseURL(conf.RedisUserURI)
	if err != nil {
		fmt.Fprint(os.Stderr, err)

		return nil
	}

	return &repositoryRedis[model.User]{
		client: redis.NewClient(opts),
		model:  &model.User{},
	}
}

func NewRfTokenRedisRepository(conf config.Config) *repositoryRedis[[]*model.RefreshToken] {
	opts, err := redis.ParseURL(conf.RedisRftURI)
	if err != nil {
		fmt.Fprint(os.Stderr, err)

		return nil
	}

	return &repositoryRedis[[]*model.RefreshToken]{
		client: redis.NewClient(opts),
		model:  &[]*model.RefreshToken{},
	}
}

func (rep *repositoryRedis[model]) Set(ctx context.Context, key string, mod model) error {
	js, err := json.Marshal(&mod)
	if err != nil {
		return fmt.Errorf("repositoryRedis, CreateUpdate, MArshal: %w", err)
	}

	_, err = rep.client.Set(ctx, key, js, time.Minute*5).Result()
	if err != nil {
		return fmt.Errorf("repositoryRedis, CreateUpdate, Set: %w", err)
	}

	return nil
}

func (rep *repositoryRedis[model]) Get(ctx context.Context, key string) (mod model, err error) {
	var js []byte
	err = rep.client.Get(ctx, key).Scan(&js)
	if err != nil {
		return mod, fmt.Errorf("repositoryRedis, GetTroughID; id nor created: %w", err)
	}

	err = json.Unmarshal(js, &mod)
	if err != nil {
		return mod, fmt.Errorf("unmarshal error at redisRepository/GetTroughID%w", err)
	}

	return mod, nil
}

func (rep *repositoryRedis[_]) Remove(ctx context.Context, key string) error {
	_, err := rep.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("repositoryRedis, failed to delete data from redise: %w", err)
	}

	return nil
}
