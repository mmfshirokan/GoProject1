package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	client *redis.Client
	user   repository.RedisRepositoryInterface[*model.User]
	token  repository.RedisRepositoryInterface[[]*model.RefreshToken]
	usrMap repository.MapRepositoryInterface[*model.User]
	tokMap repository.MapRepositoryInterface[[]*model.RefreshToken]
}

func NewConsumer(
	client *redis.Client,
	user repository.RedisRepositoryInterface[*model.User],
	token repository.RedisRepositoryInterface[[]*model.RefreshToken],
	usrMap repository.MapRepositoryInterface[*model.User],
	tokMap repository.MapRepositoryInterface[[]*model.RefreshToken],
) *Consumer {
	return &Consumer{
		client: client,
		user:   user,
		token:  token,
		usrMap: usrMap,
		tokMap: tokMap,
	}
}

func (cons *Consumer) Consume(ctx context.Context) {
	logInit()

	for {
		streams, err := repository.NewXread(ctx, cons.client)

		if err != nil {
			log.Error(fmt.Errorf("NewXRead error at consumer: %w", err))
		}
		// log.Debug(fmt.Printf("stream: %s ", streams[0])) is unnesasary - mesage is too big
		for _, message := range streams[0].Messages {
			log.Debug(fmt.Printf("ID: %s ", message.ID))

			for key, value := range message.Values {
				_, found := strings.CutPrefix(key, "user:")
				if found {
					var user model.User

					if err = json.Unmarshal([]byte(value.(string)), &user); err != nil {
						log.Error(fmt.Errorf("json.Unmarshal at consumer ReadRedisStream: %w", err))
					}
					cons.usrMap.Set(ctx, key, &user)
				}

				_, found = strings.CutPrefix(key, "token:")
				if found {
					var token []*model.RefreshToken

					if err = json.Unmarshal([]byte(value.(string)), &token); err != nil {
						log.Error(fmt.Errorf("json.Unmarshal at consumer ReadRedisStream: %w", err))
					}
					cons.tokMap.Set(ctx, key, token)
				}
			}
		}
	}
}

func logInit() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}
