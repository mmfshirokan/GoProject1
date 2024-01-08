package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepositoryInterface[object *model.User | []*model.RefreshToken] interface {
	Add(ctx context.Context, key string, obj object) error
}

type RedisRepository[object *model.User | []*model.RefreshToken] struct {
	client *redis.Client
	mut    sync.RWMutex
}

func NewCLient(conf config.Config) *redis.Client {
	opt, err := redis.ParseURL(conf.RedisURI)

	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("error ocured while parsing RedisURI (chek config.Config): %w", err))

		return nil
	}

	client := redis.NewClient(opt)
	err = client.Ping(context.Background()).Err()

	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("error ocured while conecting to redis server: %w", err))

		return nil
	}

	return client
}

func NewUserRedisRepository(client *redis.Client) RedisRepositoryInterface[*model.User] {
	return &RedisRepository[*model.User]{
		client: client,
	}
}

func NewRftRedisRepository(client *redis.Client) RedisRepositoryInterface[[]*model.RefreshToken] {
	return &RedisRepository[[]*model.RefreshToken]{
		client: client,
	}
}

func (rep *RedisRepository[object]) Add(ctx context.Context, key string, obj object) error {
	marshObj, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("marshal error ocured in repository redis: %w", err)
	}

	resObj := []string{key, string(marshObj)}

	rep.mut.Lock()

	err = rep.client.XAdd(ctx, newXAddArg(resObj)).Err()

	rep.mut.Unlock()

	if err != nil {
		return fmt.Errorf("XAdd error ocured in repository redis: %w", err)
	}

	return nil
}

func NewXread(ctx context.Context, client *redis.Client) ([]redis.XStream, error) {
	return client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"redisStr", "$"},
		Count:   2,
		Block:   0,
	}).Result()
}

func newXAddArg(obj []string) *redis.XAddArgs {
	return &redis.XAddArgs{
		Stream: "redisStr",
		ID:     strconv.FormatInt(time.Now().UnixMilli(), 10),
		Values: obj,
	}
}
