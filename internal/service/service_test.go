package service

import (
	"context"
	"errors"
	"strconv"

	"testing"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/service/mocks"
)

var (
	errSome = errors.New("someError")
)

func TestCreate(t *testing.T) {
	repoConn := mocks.NewRepositoryInterface(t)
	serviceConn := NewUser(repoConn, nil, nil)

	testTable := []struct {
		name      string
		input     model.User
		someError error
		context   context.Context
	}{
		{
			name: "standart input with ID=110",
			input: model.User{
				ID:       110,
				Name:     "Mark",
				Male:     true,
				Password: "",
			},
			someError: nil,
			context:   context.TODO(),
		},
		{
			name: "standart input with ID=113",
			input: model.User{
				ID:       113,
				Name:     "Jhon",
				Male:     true,
				Password: "",
			},
			someError: nil,
			context:   context.TODO(),
		},
		{
			name: "standart input with someError",
			input: model.User{
				ID:       1000000001,
				Name:     "Jane",
				Male:     false,
				Password: "",
			},
			someError: errSome,
			context:   context.TODO(),
		},
	}
	for _, test := range testTable {
		mockCall := repoConn.EXPECT().Create(test.context, test.input).Return(test.someError)
		serviceConn.Create(test.context, test.input)
		repoConn.AssertExpectations(t)
		mockCall.Unset()
	}
}

func TestDelete(t *testing.T) {
	repoConn := mocks.NewRepositoryInterface(t)
	mapUsrConn := mocks.NewMapRepositoryInterface[*model.User](t)
	serviceConn := NewUser(repoConn, nil, mapUsrConn)

	testTable := []struct {
		name       string
		input      int
		redisInput string
		someError  error
		context    context.Context
	}{
		{
			name:       "standart input with ID=110",
			input:      110,
			redisInput: "user:110",
			someError:  nil,
			context:    context.TODO(),
		},
		{
			name:       "standart input with ID=113",
			input:      113,
			redisInput: "user:113",
			someError:  nil,
			context:    context.TODO(),
		},
		{
			name:       "standart input with someError",
			input:      116,
			redisInput: "user:116",
			someError:  errSome,
			context:    context.TODO(),
		},
	}
	for _, test := range testTable {
		repoCall := repoConn.EXPECT().Delete(test.context, test.input).Return(test.someError)
		mapUsrCall := mapUsrConn.EXPECT().Remove(test.redisInput)
		serviceConn.Delete(test.context, test.input)
		repoConn.AssertExpectations(t)
		mapUsrConn.AssertExpectations(t)
		repoCall.Unset()
		mapUsrCall.Unset()
	}
}

func TestGetTroughID(t *testing.T) {
	repoConn := mocks.NewRepositoryInterface(t)
	redisUsrConn := mocks.NewRedisRepositoryInterface[*model.User](t)
	mapUsrConn := mocks.NewMapRepositoryInterface[*model.User](t)
	serviceConn := NewUser(
		repoConn,
		redisUsrConn,
		mapUsrConn,
	)
	testOutput := &model.User{
		ID:       110,
		Name:     "Mark",
		Male:     true,
		Password: "",
	}

	testTable := []struct {
		name    string
		input   int
		context context.Context

		usrMapOutput *model.User
		usrMapError  error

		repoOutput *model.User
		repoError  error

		redisError error
	}{
		{
			name:    "input with value stored in map",
			input:   110,
			context: context.TODO(),

			usrMapOutput: testOutput,
			usrMapError:  nil,
		},
		{
			name:    "input with value stored in repo",
			input:   113,
			context: context.TODO(),

			usrMapOutput: nil,
			usrMapError:  errSome,

			repoOutput: testOutput,
			repoError:  nil,

			redisError: nil,
		},
		{
			name:    "input with value stored in repo but GetTroughID has error",
			input:   113,
			context: context.TODO(),

			usrMapOutput: nil,
			usrMapError:  errSome,

			repoOutput: nil,
			repoError:  errSome,
		},
		{
			name:    "input with value stored in repo but redis Add has error",
			input:   113,
			context: context.TODO(),

			usrMapOutput: nil,
			usrMapError:  errSome,

			repoOutput: testOutput,
			repoError:  nil,

			redisError: errSome,
		},
	}
	for _, test := range testTable {
		inputKey := "user:" + strconv.FormatInt(int64(test.input), 10)
		mapCall := mapUsrConn.EXPECT().Get(test.context, inputKey).Return(test.usrMapOutput, test.usrMapError)
		repoCall := repoConn.EXPECT().GetTroughID(test.context, test.input).Return(test.repoOutput, test.repoError).Maybe()
		redisCall := redisUsrConn.EXPECT().Add(test.context, inputKey, test.repoOutput).Return(test.redisError).Maybe()

		serviceConn.GetTroughID(test.context, test.input)

		mapUsrConn.AssertExpectations(t)
		repoConn.AssertExpectations(t)
		redisUsrConn.AssertExpectations(t)

		if test.usrMapError != nil {
			repoConn.AssertCalled(t, "GetTroughID", test.context, test.input)
		}
		if test.usrMapError != nil && test.repoError != nil {
			redisUsrConn.AssertNotCalled(t, "Add", test.context, inputKey, test.repoOutput)
		}

		mapCall.Unset()
		redisCall.Unset()
		repoCall.Unset()
	}
}
