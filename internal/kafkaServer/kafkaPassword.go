package kafkaserver

import (
	"context"
	"io"
	"strconv"
	"sync"

	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/segmentio/kafka-go"
)

type KafkaUser struct {
	password repository.PwRepositoryInterface
	mut      sync.Mutex
}

type KafkaUserInterface interface {
	Read(context context.Context, broker string, topic string)
}

func NewKafkaServer(pwd repository.PwRepositoryInterface) KafkaUserInterface {
	return &KafkaUser{
		password: pwd,
	}
}

func (pwd *KafkaUser) Read(ctx context.Context, broker string, topic string) {
	for i := 0; i <= 2; i++ {
		pwd.mut.Lock()
		go func() {
			reader := kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{broker},
				Topic:     topic,
				Partition: i - 1,
			})
			pwd.mut.Unlock()
			defer reader.Close()

			numOfRows := 0
			resArr := make([][]interface{}, 100)
			//gorutineNum := reader.Config().Partition

			//log.Info("Sucssesful start ", gorutineNum)

			for {
				msg, err := reader.ReadMessage(ctx)
				if err == io.EOF {
					//log.Info(fmt.Sprintf("exiting kafka read gorutine-%d", gorutineNum))
					break
				}
				if err != nil {
					//log.Error(fmt.Sprintf("error occured in gorutine-%d: %v", gorutineNum, err))
					break
				}

				id, err := strconv.ParseInt(string(msg.Key), 10, 64)
				if err != nil {
					//log.Error(fmt.Sprintf("error occured in gorutine-%d: %v", gorutineNum, err))
					continue
				}

				resArr[numOfRows] = []interface{}{int(id), string(msg.Value)}
				numOfRows++

				if numOfRows > 99 {
					pwd.mut.Lock()
					err := pwd.password.BulkStore(ctx, resArr)
					pwd.mut.Unlock()

					if err != nil {
						//log.Error(fmt.Printf("error occured in gorutine-%d: %v", gorutineNum, err))
					}
					numOfRows = 0
				}
			}
		}()
	}
}
