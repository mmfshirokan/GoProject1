package service

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	uuidArr = [3]uuid.UUID{
		uuid.MustParse("a41222ff-3750-463c-82de-3850836afd7e"),
		uuid.MustParse("bbe8b692-0a16-49fa-9700-355842b3f6e0"),
		uuid.MustParse("1d89d00e-ee9f-4ab5-9dd9-270e5d800a46"),
	}
)

func TestCreateAuthToken(t *testing.T) { // IS it nessasary for create token to be *Token method
	authServiceConn := NewToken(nil, nil, nil)
	testTable := []struct {
		name      string
		inputID   int
		inputName string
		inputMale bool
	}{
		{
			name:      "std with ID=110",
			inputID:   110,
			inputName: "Jhon",
			inputMale: true,
		},
		{
			name:      "std input with ID=113",
			inputID:   113,
			inputName: "Jane",
			inputMale: false,
		},
	}
	for _, test := range testTable {
		authToken := authServiceConn.CreateAuthToken(test.inputID, test.inputName, test.inputMale)
		parsedAuthToken, err := jwt.ParseWithClaims(authToken, &model.UserRequest{}, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			t.Error("token parsing failed therefore methodis wrong")
		}

		payload, ok := parsedAuthToken.Claims.(*model.UserRequest)
		if !ok {
			t.Error("Wrong interface assertion in TestCreateAuthToken")
		}
		assert.Equal(t, test.inputID, payload.ID)
		assert.Equal(t, test.inputMale, payload.Male)
		assert.Equal(t, test.inputName, payload.Name)
	}
}

func TestCreateRfToken(t *testing.T) {
	authRepoConn := mocks.NewAuthRepositoryInterface(t)
	authServiceConn := NewToken(authRepoConn, nil, nil)

	testTable := []struct {
		name      string
		input     int
		context   context.Context
		someError error
	}{
		{
			name:      "std input with ID=110",
			input:     110,
			context:   context.TODO(),
			someError: nil,
		},
		{
			name:      "std input with ID=113",
			input:     113,
			context:   context.TODO(),
			someError: nil,
		},
		{
			name:      "input with someError",
			input:     133,
			context:   context.TODO(),
			someError: errSome,
		},
	}
	for _, test := range testTable {
		authCall := authRepoConn.EXPECT().Create(test.context, mock.Anything).Return(test.someError)
		authServiceConn.CreateRfToken(test.context, test.input)
		authRepoConn.AssertExpectations(t)
		authCall.Unset()
	}
}

func TestValidateRfTokenTroughID(t *testing.T) {
	authServiceConn := NewToken(nil, nil, nil)
	var (
		hashArr [3]string
		err     error
	)

	for i, val := range uuidArr {
		if hashArr[i], err = conductHashing(val); err != nil {
			t.Error("conductHashing error")
		}
	}

	testTable := []struct {
		name      string
		context   context.Context
		inputUUID uuid.UUID
		inputHash string
		output    bool
	}{
		{
			name:      "std input with correct hash",
			context:   context.TODO(),
			inputUUID: uuidArr[0],
			inputHash: hashArr[0],
			output:    true,
		},
		{
			name:      "std input with incorrect hash",
			context:   context.TODO(),
			inputUUID: uuidArr[1],
			inputHash: hashArr[0],
			output:    false,
		},
		{
			name:      "second std input with correct hash",
			context:   context.TODO(),
			inputUUID: uuidArr[2],
			inputHash: hashArr[2],
			output:    true,
		},
	}
	for _, test := range testTable {
		recivedOutput, _ := authServiceConn.ValidateRfTokenTroughID(test.inputHash, test.inputUUID)
		assert.Equal(t, test.output, recivedOutput)
	}
}

func TestAuthDelete(t *testing.T) {
	authRepoConn := mocks.NewAuthRepositoryInterface(t)
	authServiceConn := NewToken(authRepoConn, nil, nil)
	testTable := []struct {
		name      string
		input     uuid.UUID
		context   context.Context
		someError error
	}{
		{
			name:      "std input with ID=110",
			input:     uuidArr[0],
			context:   context.TODO(),
			someError: nil,
		},
		{
			name:      "std input with ID=113",
			input:     uuidArr[1],
			context:   context.TODO(),
			someError: nil,
		},
		{
			name:      "input with someError",
			input:     uuidArr[2],
			context:   context.TODO(),
			someError: errSome,
		},
	}
	for _, test := range testTable {
		authRepoCall := authRepoConn.EXPECT().Delete(test.context, test.input).Return(test.someError)
		authServiceConn.Delete(test.context, test.input)
		authRepoCall.Unset()
	}
}

func TestGetByUserID(t *testing.T) {
	authRepoConn := mocks.NewAuthRepositoryInterface(t)
	tokenRedisConn := mocks.NewRedisRepositoryInterface[[]*model.RefreshToken](t)
	tokenMapConn := mocks.NewMapRepositoryInterface[[]*model.RefreshToken](t)
	authServiceConn := NewToken(authRepoConn, tokenRedisConn, tokenMapConn)
	refreshLifeTime := time.Now().Add(time.Hour * 12)

	testTable := []struct {
		name    string
		input   int
		context context.Context

		mapOutput []*model.RefreshToken
		mapError  error

		repoOutput []*model.RefreshToken
		repoError  error

		reidsError error
	}{
		{
			name:    "std input with ID=110",
			input:   110,
			context: context.TODO(),
			mapOutput: []*model.RefreshToken{
				{
					UserID:     110,
					ID:         uuidArr[0],
					Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
					Expiration: refreshLifeTime,
				},
				{
					UserID:     110,
					ID:         uuidArr[1],
					Hash:       "Y230k558gdR_zR-QR9dvesHkO5gsaZE4dsX2yahJdU0",
					Expiration: refreshLifeTime,
				},
			},
			mapError: nil,
		},
		{
			name:       "std input with ID=113, mapError and repoError",
			input:      113,
			context:    context.TODO(),
			mapOutput:  nil,
			mapError:   errSome,
			repoOutput: nil,
			repoError:  errSome,
		},
		{
			name:      "std input with ID=113 and mapError",
			input:     113,
			context:   context.TODO(),
			mapOutput: nil,
			mapError:  errSome,
			repoOutput: []*model.RefreshToken{
				{
					UserID:     113,
					ID:         uuidArr[2],
					Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
					Expiration: refreshLifeTime,
				},
			},
			repoError:  nil,
			reidsError: nil,
		},
		{
			name:      "std input with ID=113, mapError and redisError",
			input:     113,
			context:   context.TODO(),
			mapOutput: nil,
			mapError:  errSome,
			repoOutput: []*model.RefreshToken{
				{
					UserID:     113,
					ID:         uuidArr[2],
					Hash:       "Qxm0k58B14G_zR-QR9dvesHkO56yKZE48rX2yahJdU0",
					Expiration: refreshLifeTime,
				},
			},
			repoError:  nil,
			reidsError: errSome,
		},
	}
	for _, test := range testTable {
		inputKey := "token:" + strconv.FormatInt(int64(test.input), 10)
		tokenMapCall := tokenMapConn.EXPECT().Get(test.context, inputKey).Return(test.mapOutput, test.mapError)
		authRepoCall := authRepoConn.EXPECT().GetByUserID(test.context, test.input).Return(test.repoOutput, test.repoError).Maybe()
		tokenRedisCall := tokenRedisConn.EXPECT().Add(test.context, inputKey, test.repoOutput).Return(test.reidsError).Maybe()

		authServiceConn.GetByUserID(test.context, test.input)

		tokenMapConn.AssertExpectations(t)
		authRepoConn.AssertExpectations(t)
		tokenRedisConn.AssertExpectations(t)

		if test.mapError != nil {
			authRepoConn.AssertCalled(t, "GetByUserID", test.context, test.input)
		}
		if test.mapError != nil && test.repoError != nil {
			tokenRedisConn.AssertNotCalled(t, "Add", test.context, inputKey, test.repoOutput)
		}

		tokenMapCall.Unset()
		authRepoCall.Unset()
		tokenRedisCall.Unset()
	}
}
