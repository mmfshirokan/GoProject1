package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/stretchr/testify/assert"
)

const numberOfTestCases int = 4

var (
	authConn        AuthRepositoryInterface
	uuidsArr        [numberOfTestCases]uuid.UUID
	refreshLifeTime time.Time
)

func TestAuthCreate(t *testing.T) {
	type testCase struct {
		name     string
		input    *model.RefreshToken
		hasError bool
	}

	for i := 0; i < len(uuidsArr); i++ {
		uuidsArr[i] = uuid.New()
	}

	refreshLifeTime := time.Now().Add(time.Hour * 12)

	testTable := []testCase{
		{
			name: "standart input with ID=110",
			input: &model.RefreshToken{
				UserID:     110,
				ID:         uuidsArr[0],
				Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
				Expiration: refreshLifeTime,
			},
			hasError: false,
		},
		{
			name: "standart input with ID=110 and duplicate UUID",
			input: &model.RefreshToken{
				UserID:     110,
				ID:         uuidsArr[0],
				Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
				Expiration: refreshLifeTime,
			},
			hasError: true,
		},
		{
			name: "second standart input with ID=110", // TODO add the same UUID
			input: &model.RefreshToken{
				UserID:     110,
				ID:         uuidsArr[1],
				Hash:       "Y230k558gdR_zR-QR9dvesHkO5gsaZE4dsX2yahJdU0",
				Expiration: refreshLifeTime,
			},
			hasError: false,
		},
		{
			name: "standart input with ID=113",
			input: &model.RefreshToken{
				UserID:     113,
				ID:         uuidsArr[2],
				Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
				Expiration: refreshLifeTime,
			},
			hasError: false,
		},
		{
			name: "standart input with ID bigger than allowed",
			input: &model.RefreshToken{
				UserID:     1000000001,
				ID:         uuidsArr[3],
				Hash:       "Qxm0k58B14G_zR-QR9dvesH312312312",
				Expiration: refreshLifeTime,
			},
			hasError: true,
		},
	}
	for _, test := range testTable {
		err := authConn.Create(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Print("TestAuthCreate Finished!\n")
}

func TestGetByUserID(t *testing.T) {
	type testCase struct {
		name     string
		input    int
		output   []*model.RefreshToken
		hasError bool
	}
	//sd := make([]*model.RefreshToken, 2)
	testTable := []testCase{
		{
			name:  "standart input with ID=110",
			input: 110,
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
		},
		{
			name:  "standart input with ID=113",
			input: 113,
			output: []*model.RefreshToken{
				{
					UserID:     113,
					ID:         uuidsArr[2],
					Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
					Expiration: refreshLifeTime,
				},
			},
			hasError: false,
		},
		{
			name:     "standart input with ID that do not exist",
			input:    1334,
			output:   nil,
			hasError: false,
		},
	}
	for _, test := range testTable {
		actuallOutput, err := authConn.GetByUserID(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, test.name)
			assert.Nil(t, test.output, "otput must be nil: ", test.name)
		} else {
			assert.Nil(t, err, test.name)
			//assert.Equal(t, test.output, actuallOutput) //add time comparison
			for i := 0; i < (len(actuallOutput) - 2); i++ {
				assert.Equal(t, test.output[i], actuallOutput[i])
			}
		}
	}
	fmt.Print("TestGetByID Finished!\n")
}

func TestAuthDelete(t *testing.T) {
	type testCase struct {
		name     string
		input    uuid.UUID
		hasError bool
	}
	testTable := []testCase{
		{
			name:     "standart input with ID=110",
			input:    uuidsArr[0],
			hasError: false,
		},
		{
			name:     "second standart input with ID=110 and different uuid",
			input:    uuidsArr[1],
			hasError: false,
		},
		{
			name:     "standart input with ID=113",
			input:    uuidsArr[2],
			hasError: false,
		},
		{
			name:     "standart input with UUID that do not exist", // TODO cahnge for Error
			input:    uuidsArr[3],
			hasError: false,
		},
	}
	for _, test := range testTable {
		err := authConn.Delete(context.Background(), test.input)
		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
	fmt.Print("TestAuthDelete Finished!\n")
}
