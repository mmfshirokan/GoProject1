package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/stretchr/testify/assert"
)

var (
	rfTokenExpirationTime time.Time
	redisUsrConn          *repository.RedisRepository[*model.User]
	redisRftConn          *repository.RedisRepository[[]*model.RefreshToken]
)

func TestSet(t *testing.T) {
	type testCase[obj *model.User | []*model.RefreshToken] struct {
		name     string
		inputKey string
		inputObl obj
		hasError bool
	}
	usrTestTable := []testCase[*model.User]{
		{
			name:     "standart input with ID=110",
			inputKey: "id:110",
			inputObl: &model.User{
				ID:       110,
				Name:     "Mark",
				Male:     true,
				Password: "abcd",
			},
			hasError: false,
		},
		{
			name: "standart input with ID=113",
			inputObl: &model.User{
				ID:       113,
				Name:     "Jhane",
				Male:     false,
				Password: "s0pranO",
			},
			hasError: false,
		},
	}

	rfTokenExpirationTime = time.Now().Add(time.Hour * 6)

	rftTestTable := []testCase[[]*model.RefreshToken]{
		{
			name:     "standart input with ID=110",
			inputKey: "id:110",
			inputObl: []*model.RefreshToken{
				{
					UserID:     110,
					ID:         uuid.New(),
					Hash:       "aboba",
					Expiration: rfTokenExpirationTime,
				},
			},
			hasError: false,
		},
		{
			name: "standart input with ID=113",
			inputObl: []*model.RefreshToken{
				{
					UserID:     113,
					ID:         uuid.New(),
					Hash:       "aboba",
					Expiration: rfTokenExpirationTime,
				},
				{
					UserID:     113,
					ID:         uuid.New(),
					Hash:       "abiba",
					Expiration: rfTokenExpirationTime,
				},
			},
			hasError: false,
		},
	}

	for _, test := range usrTestTable {
		err := redisUsrConn.Add(context.Background(), test.inputKey, test.inputObl)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}

	for _, test := range rftTestTable {
		err := redisRftConn.Add(context.Background(), test.inputKey, test.inputObl)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Println("TestSet Finished!")
}
