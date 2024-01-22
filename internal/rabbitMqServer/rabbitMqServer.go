package rabbitmqserver

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/mmfshirokan/GoProject1/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMqUser struct {
	password repository.PwRepositoryInterface
	mut      sync.Mutex
}

type RabbitMqInterface interface {
	Read(ctx context.Context, url string, keys []string)
}

func NewRabbitMqServer(pwd repository.PwRepositoryInterface) RabbitMqInterface {
	return &RabbitMqUser{
		password: pwd,
	}
}

func (rbmq *RabbitMqUser) Read(ctx context.Context, url string, keys []string) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("connection error: %v", err)
		return
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("chanel opening error: %v", err)
		return
	}
	//defer ch.Close()

	for i, key := range keys {
		go func(key string) {
			q, err := ch.QueueDeclare(
				key,   // name
				false, // durable
				false, // delete when unused
				false, // exclusive
				false, // no-wait
				nil,   // arguments
			)
			if err != nil {
				log.Errorf("Queue-declare %v error occurred: %v", q.Name, err)
				return
			}

			messages, err := ch.Consume(
				q.Name, // queue
				"",     // consumer
				true,   // auto-ack
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // args
			)
			if err != nil {
				log.Errorf("Consume error occurred: %v", err)
				return
			}

			var forever chan struct{}

			go func() {
				numOfRows := 0
				resArr := make([][]interface{}, 100)

				for message := range messages {
					id, err := strconv.ParseInt(string(message.Body), 10, 64)
					if err != nil {
						log.Error(fmt.Sprintf("error occured in gorutine-%d: %v", i, err))
						continue
					}

					resArr[numOfRows] = []interface{}{int(id), "pwd" + string(message.Body)}
					numOfRows++

					if numOfRows > 99 {
						rbmq.mut.Lock()
						err := rbmq.password.BulkStore(ctx, resArr)
						rbmq.mut.Unlock()

						if err != nil {
							log.Error(fmt.Printf("error occured in gorutine-%d: %v", i, err))
						}
						numOfRows = 0
					}
				}
			}()

			<-forever

		}(key)
	}
}
