package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/stretchr/testify/assert"
)

var pwConn repository.PwRepositoryInterface

func TestStore(t *testing.T) {
	type testCase struct {
		name     string
		input    model.User
		hasError bool
	}
	testTable := []testCase{
		{
			name: "standart input with ID=110",
			input: model.User{
				ID:       110,
				Password: "abcd",
			},
			hasError: false,
		},
		{
			name: "standart input with ID=113",
			input: model.User{
				ID:       113,
				Password: "s0pranO",
			},
			hasError: false,
		},
		{
			name: "password biger than allowed",
			input: model.User{
				ID:       114,
				Password: "12345678901234567890123456789012345678901",
			},
			hasError: true,
		},
	}

	for _, test := range testTable {
		err := pwConn.Store(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Println("TestStore Finished!")
}

func TestCompare(t *testing.T) {
	type testCase struct {
		name     string
		input    model.User
		output   bool
		hasError bool
	}
	testTable := []testCase{
		{
			name: "standart input with id=110",
			input: model.User{
				ID:       110,
				Password: "abcd",
			},
			output:   true,
			hasError: false,
		},
		{
			name: "standart input with id=113",
			input: model.User{
				ID:       113,
				Password: "s0pranO",
			},
			output:   true,
			hasError: false,
		},
		{
			name: "wrong password with id=113",
			input: model.User{
				ID:       113,
				Password: "Aboba167",
			},
			output:   false,
			hasError: false,
		},
		{
			name: "id that does not exists",
			input: model.User{
				ID:       41214,
				Password: "asd",
			},
			output:   false,
			hasError: true,
		},
	}

	for _, test := range testTable {
		passBool, err := pwConn.Compare(context.Background(), test.input)
		assert.Equal(t, test.output, passBool)

		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Println("TestCompare Finished!")
}

func TestDeletePassword(t *testing.T) {
	type testCase struct {
		name     string
		input    int
		hasError bool
	}
	testTable := []testCase{
		{
			name:     "standart input with ID=110",
			input:    110,
			hasError: false,
		},
		{
			name:     "standart input with ID=113",
			input:    113,
			hasError: false,
		},
		{ //no error
			name:     "ID that does not exists",
			input:    1121,
			hasError: false,
		},
	}

	for _, test := range testTable {
		err := pwConn.DeletePassword(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Println("TestDeletePassword Finished!")
}
