package handlers

import (
	"errors"

	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/internal/handlers/mocks"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	errSome    error     = errors.New("someError")
	jsonUsrArr [3]string = [3]string{
		`{"id":110,"name":"Jhon","male":true,"password":"abcd"}`,
		`{"id":113,"name":"Jane","male":false,"password":"s0pranO"}`,
		`{"id":133,"name":"Joanna","male":false,"password":"7788"}`,
	}

	userArr [3]model.User = [3]model.User{
		{
			ID:       110,
			Name:     "Jhon",
			Male:     true,
			Password: "abcd",
		},
		{
			ID:       113,
			Name:     "Jane",
			Male:     false,
			Password: "s0pranO",
		},
		{
			ID:       133,
			Name:     "Joanna",
			Male:     false,
			Password: "7788",
		},
	}
	rfTokenExpirationTime                       = time.Now().Add(time.Hour * 6)
	testingRfTokens       []*model.RefreshToken = []*model.RefreshToken{
		{
			UserID:     110,
			ID:         uuid.New(),
			Hash:       "aboba",
			Expiration: rfTokenExpirationTime,
		},
		{
			UserID:     110,
			ID:         uuid.New(),
			Hash:       "abiba",
			Expiration: rfTokenExpirationTime,
		},
	}
)

func TestSingUp(t *testing.T) {
	user := mocks.NewUserInterface(t)
	password := mocks.NewPasswordInterface(t)
	handler := NewHandler(user, password, nil)

	handlerTarget := "/users/signup"
	handlerMethod := http.MethodPost

	e := echo.New()
	rec := httptest.NewRecorder()

	testTable := []struct {
		name        string
		body        string
		userInput   model.User
		repoError   error
		pwrepoError error
	}{
		{
			name: "std input with ID=110",
			body: jsonUsrArr[0],
			userInput: model.User{
				ID:       110,
				Name:     "Jhon",
				Male:     true,
				Password: "abcd",
			},
			repoError:   nil,
			pwrepoError: nil,
		},
		{
			name: "std input with ID=113",
			body: jsonUsrArr[1],
			userInput: model.User{
				ID:       113,
				Name:     "Jane",
				Male:     false,
				Password: "s0pranO",
			},
			repoError:   nil,
			pwrepoError: nil,
		},
		{
			name: "input with repo error",
			body: jsonUsrArr[2],
			userInput: model.User{
				ID:       133,
				Name:     "Joanna",
				Male:     false,
				Password: "7788",
			},
			repoError:   errSome,
			pwrepoError: nil,
		},
		{
			name: "input with pwRepo error",
			body: jsonUsrArr[2],
			userInput: model.User{
				ID:       133,
				Name:     "Joanna",
				Male:     false,
				Password: "7788",
			},
			repoError:   nil,
			pwrepoError: errSome,
		},
	}
	for _, test := range testTable {
		req := httptest.NewRequest(handlerMethod, handlerTarget, strings.NewReader(test.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		ctx := c.Request().Context()

		usrCall := user.EXPECT().Create(ctx, test.userInput).Return(test.repoError)
		pwCall := password.EXPECT().Store(ctx, test.userInput).Return(test.pwrepoError).Maybe()

		err := handler.SignUp(c)

		user.AssertExpectations(t)
		password.AssertExpectations(t)

		if test.repoError != nil {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, test.repoError.Error()), err)
			password.AssertNotCalled(t, "Store", ctx, test.userInput)
		} else if test.pwrepoError != nil {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, test.pwrepoError.Error()), err)
		} else {
			assert.Nil(t, err, test.name)
		}

		usrCall.Unset()
		pwCall.Unset()
	}
}

func TestSignIn(t *testing.T) {
	password := mocks.NewPasswordInterface(t)
	token := mocks.NewTokenInterface(t)
	handler := NewHandler(nil, password, token)

	handlerTarget := "/users/signin"
	handlerMethod := http.MethodPut

	e := echo.New()
	rec := httptest.NewRecorder()

	testTable := []struct {
		name                 string
		body                 string
		passwordCorrect      bool
		passwordCompareError error
		userModel            model.User
		createRfTokenError   error
		getByUserIdError     error
	}{
		{
			name:                 "std input with ID=110 and incorrect password",
			body:                 jsonUsrArr[0],
			passwordCorrect:      false,
			passwordCompareError: nil,
			userModel:            userArr[0],
		},
		{
			name:                 "std input with ID=110 and correct password",
			body:                 jsonUsrArr[0],
			passwordCorrect:      true,
			passwordCompareError: nil,
			userModel:            userArr[0],
			createRfTokenError:   nil,
			getByUserIdError:     nil,
		},
		{
			name:                 "std input with ID=113 and correct password",
			body:                 jsonUsrArr[1],
			passwordCorrect:      true,
			passwordCompareError: nil,
			userModel:            userArr[1],
			createRfTokenError:   nil,
			getByUserIdError:     nil,
		},
		{
			name:                 "input with passwordCompareError and ID=133",
			body:                 jsonUsrArr[2],
			passwordCorrect:      true,
			userModel:            userArr[2],
			passwordCompareError: errSome,
		},
		{
			name:                 "input with createRfTokenError and ID=133",
			body:                 jsonUsrArr[2],
			passwordCorrect:      true,
			passwordCompareError: nil,
			userModel:            userArr[2],
			createRfTokenError:   errSome,
		},
		{
			name:                 "input with GetBYUserIdError and ID=133",
			body:                 jsonUsrArr[2],
			passwordCorrect:      true,
			passwordCompareError: nil,
			userModel:            userArr[2],
			createRfTokenError:   nil,
			getByUserIdError:     errSome,
		},
	}
	for _, test := range testTable {
		req := httptest.NewRequest(handlerMethod, handlerTarget, strings.NewReader(test.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		ctx := c.Request().Context()

		pwCall := password.EXPECT().Compare(ctx, test.userModel).Return(test.passwordCorrect, test.passwordCompareError)
		crtAuthTokenCall := token.EXPECT().CreateAuthToken(test.userModel.ID, test.userModel.Name, test.userModel.Male).Return(mock.Anything).Maybe()
		crtRfTokenCall := token.EXPECT().CreateRfToken(ctx, test.userModel.ID).Return(test.createRfTokenError).Maybe()
		getRfToken := token.EXPECT().GetByUserID(ctx, test.userModel.ID).Return(testingRfTokens, test.getByUserIdError).Maybe()

		handler.SignIn(c)

		password.AssertExpectations(t)
		token.AssertExpectations(t)

		pwCall.Unset()
		crtAuthTokenCall.Unset()
		crtRfTokenCall.Unset()
		getRfToken.Unset()
	}
}

func TestRefresh(t *testing.T) { // add more tests
	user := mocks.NewUserInterface(t)
	token := mocks.NewTokenInterface(t)
	handler := NewHandler(user, nil, token)

	testingUUID := "3e0fa715-3e4f-4568-8e81-2a0ba8ea575b"
	testingHash := "AqSmYqlyqhYJezmvx-aWtbXAzj-0MYiHBsycNc81nck"
	experationTime := time.Now().Add(time.Hour * 6)

	handlerTarget := "/users/refresh"
	handlerMethod := http.MethodPut

	e := echo.New()
	rec := httptest.NewRecorder()

	testTable := []struct {
		name string
		body string
		id   int

		tokenValide        bool
		tokenValidateError error
		tokenDeleteError   error
		createRfTokenError error

		getByUserIdOutput []*model.RefreshToken
		getByUserIdError  error

		getTroughIdOutput    *model.User
		getTroughIdError     error
		creataAuthTokenError error
	}{
		{
			name: "std input with ID=110",
			body: fmt.Sprintf(`"userID":"110","id":"%v","hash":"%s","expiration":%v`, testingUUID, testingHash, experationTime),
			id:   110,

			tokenValide:        true,
			tokenValidateError: nil,
			tokenDeleteError:   nil,
			createRfTokenError: nil,

			getByUserIdOutput: []*model.RefreshToken{
				{
					UserID:     110,
					ID:         uuid.New(),
					Hash:       "someHash",
					Expiration: experationTime,
				},
			},
			getByUserIdError: nil,

			getTroughIdOutput:    &userArr[0],
			getTroughIdError:     nil,
			creataAuthTokenError: nil,
		},
	}

	for _, test := range testTable {
		req := httptest.NewRequest(handlerMethod, handlerTarget, strings.NewReader(test.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		ctx := c.Request().Context()

		validateCall := token.EXPECT().ValidateRfTokenTroughID(testingHash, uuid.MustParse(testingUUID)).Return(test.tokenValide, test.tokenValidateError)
		deleteCall := token.EXPECT().Delete(ctx, uuid.MustParse(testingUUID)).Return(test.tokenDeleteError)
		rftCreateCall := token.EXPECT().CreateRfToken(ctx, test.id).Return(test.createRfTokenError).Maybe()
		getByUserCall := token.EXPECT().GetByUserID(ctx, test.id).Return(test.getByUserIdOutput, test.getByUserIdError).Maybe()
		getTroughIdCall := user.EXPECT().GetTroughID(ctx, test.id).Return(test.getTroughIdOutput, test.getTroughIdError).Maybe()
		createAuthCall := token.EXPECT().CreateAuthToken(test.getTroughIdOutput.ID, test.getTroughIdOutput.Name, test.getTroughIdOutput.Male).Return(mock.Anything).Maybe()

		handler.Refresh(c)

		validateCall.Unset()
		deleteCall.Unset()
		rftCreateCall.Unset()
		getByUserCall.Unset()
		getTroughIdCall.Unset()
		createAuthCall.Unset()

	}
}
