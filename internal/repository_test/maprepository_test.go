package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/stretchr/testify/assert"
)

var (
	mapUsrConn *repository.MapRepository[*model.User]
	mapRftConn *repository.MapRepository[[]*model.RefreshToken]
	mapUsr     map[string]*model.User
	mapRft     map[string][]*model.RefreshToken
)

func TestMapSet(t *testing.T) {
	type testCase[obj *model.User | []*model.RefreshToken] struct {
		name     string
		inputKey string
		inputObj obj
		context  context.Context
	}
	usrTestTable := []testCase[*model.User]{
		{
			name:     "standart input with inputKey=user:110",
			inputKey: "user:110",
			inputObj: &model.User{
				ID:       110,
				Name:     "Jhon",
				Male:     true,
				Password: "abcd",
			},
			context: context.TODO(),
		},
		{
			name:     "standart input witg inputKey=user:113",
			inputKey: "user:113",
			inputObj: &model.User{
				ID:       113,
				Name:     "Jane",
				Male:     false,
				Password: "s0pranO",
			},
			context: context.TODO(),
		},
	}
	for _, test := range usrTestTable {
		mapUsrConn.Set(test.context, test.inputKey, test.inputObj)

		_, ok := mapUsr[test.inputKey]
		if !ok {
			t.Error("TestMapSet failed")
		}
	}

	rftTestTable := []testCase[[]*model.RefreshToken]{
		{
			name:     "standart input withe inputKey=token:110",
			inputKey: "token:110",
			inputObj: []*model.RefreshToken{
				{
					UserID:     110,
					ID:         uuidsArr[0],
					Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
					Expiration: refreshLifeTime,
				},
				{
					UserID:     110,
					ID:         uuidsArr[1],
					Hash:       "Y230k558gdR_zR-QR9dvesHkO5gsaZE4dsX2yahJdU0",
					Expiration: refreshLifeTime,
				},
			},
			context: context.TODO(),
		},
		{
			name:     "standart input witg inputKey=token:113",
			inputKey: "token:113",
			inputObj: []*model.RefreshToken{
				{
					UserID:     113,
					ID:         uuidsArr[2],
					Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
					Expiration: refreshLifeTime,
				},
			},
			context: context.TODO(),
		},
	}
	for _, test := range rftTestTable {
		mapRftConn.Set(test.context, test.inputKey, test.inputObj)

		_, ok := mapRft[test.inputKey]
		if !ok {
			t.Error("TestMapSet failed with: " + test.name)
		}
	}
	fmt.Println("TestMapSet Finished!")
}

func TestGet(t *testing.T) {
	type testCase[object *model.User | []*model.RefreshToken] struct {
		name     string
		inputKey string
		output   object
		context  context.Context
		hasError bool
	}
	usrTestTable := []testCase[*model.User]{
		{
			name:     "standart input with inputKey=user:110",
			inputKey: "user:110",
			output: &model.User{
				ID:       110,
				Name:     "Jhon",
				Male:     true,
				Password: "abcd",
			},
			hasError: false,
			context:  context.TODO(),
		},
		{
			name:     "standart input witg inputKey=user:113",
			inputKey: "user:113",
			output: &model.User{
				ID:       113,
				Name:     "Jane",
				Male:     false,
				Password: "s0pranO",
			},
			hasError: false,
			context:  context.TODO(),
		},
	}
	for _, test := range usrTestTable {
		res, err := mapUsrConn.Get(test.context, test.inputKey)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Equal(t, test.output, res)
		}
	}

	rftTestTable := []testCase[[]*model.RefreshToken]{
		{
			name:     "standart input with inputKey=token:110",
			inputKey: "token:110",
			output: []*model.RefreshToken{
				{
					UserID:     110,
					ID:         uuidsArr[0],
					Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
					Expiration: refreshLifeTime,
				},
				{
					UserID:     110,
					ID:         uuidsArr[1],
					Hash:       "Y230k558gdR_zR-QR9dvesHkO5gsaZE4dsX2yahJdU0",
					Expiration: refreshLifeTime,
				},
			},
			hasError: false,
			context:  context.TODO(),
		},
		{
			name:     "standart input witg inputKey=token:113",
			inputKey: "token:113",
			output: []*model.RefreshToken{
				{
					UserID:     113,
					ID:         uuidsArr[2],
					Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
					Expiration: refreshLifeTime,
				},
			},
			hasError: false,
			context:  context.TODO(),
		},
	}
	for _, test := range rftTestTable {
		res, err := mapRftConn.Get(test.context, test.inputKey)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Equal(t, test.output, res)
		}
	}
	fmt.Println("TestGet Finished!")
}

func TestRemove(t *testing.T) { // add concurent testing
	type testCase struct {
		name      string
		inputKey  string
		isUserMap bool
	}
	usrTestTable := []testCase{
		{
			name:      "standart input with inputKey=user:110",
			inputKey:  "user:110",
			isUserMap: true,
		},
		{
			name:      "standart input with inputKey=user:113",
			inputKey:  "user:113",
			isUserMap: true,
		},
		{
			name:      "standart input with inputKey=token:110",
			inputKey:  "token:110",
			isUserMap: false,
		},
		{
			name:      "standart input with inputKey=token:113",
			inputKey:  "token:113",
			isUserMap: false,
		},
	}
	for _, test := range usrTestTable {
		if test.isUserMap {
			mapUsrConn.Remove(test.inputKey)
			_, ok := mapUsr[test.inputKey]
			if ok {
				t.Error("item wasn't removed in TestRemove with: " + test.inputKey)
			}
		} else {
			mapRftConn.Remove(test.inputKey)
			_, ok := mapRft[test.inputKey]
			if ok {
				t.Error("item wasn't removed in TestRemove with: " + test.inputKey)
			}
		}
	}
	fmt.Println("TestRemove Finished!")
}
