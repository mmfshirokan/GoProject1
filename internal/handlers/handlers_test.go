package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/mmfshirokan/GoProject1/internal/handlers/mocks"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/stretchr/testify/assert"
)

var (
	testingJwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, &model.UserRequest{
		ID:   110,
		Name: "Jhon",
		Male: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(rfTokenExpirationTime),
		},
	})
	testingModelUser = model.User{
		ID:   110,
		Name: "Jhon",
		Male: true,
	}
)

func TestGetUser(t *testing.T) { // add more test
	handler := NewHandler(nil, nil, nil)
	testingMethod := http.MethodGet
	testingTarget := "/users/auth/get"

	req := httptest.NewRequest(testingMethod, testingTarget, nil)
	rec := httptest.NewRecorder()
	e := echo.New()

	testTable := []struct {
		name     string
		jwtToken *jwt.Token
		hasError bool

		userId   string
		userName string
		userMale string
	}{
		{
			name:     "std input with ID=110",
			jwtToken: testingJwtToken,
			hasError: false,

			userId:   "110",
			userName: "Jhon",
			userMale: "true",
		},
		{
			name:     "std input with ID=110",
			jwtToken: jwt.New(jwt.SigningMethodHS256),
			hasError: true,
		},
	}

	for _, test := range testTable {
		c := e.NewContext(req, rec)
		c.Set("user", test.jwtToken)
		err := handler.GetUser(c)
		bodyString := fmt.Sprint(
			"Usser id: ",
			test.userId,
			"\nUser name: ",
			test.userName,
			"\nUser male:",
			test.userMale,
			"\n",
		)

		if test.hasError {
			assert.Error(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, bodyString, rec.Body.String())
		}
	}
}

func TestUpdateUser(t *testing.T) {
	user := mocks.NewUserInterface(t)
	handler := NewHandler(user, nil, nil)

	method := http.MethodPut
	target := "/users/auth/update"

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, nil)

	testTable := []struct {
		name     string
		jwtToken *jwt.Token

		updtaeInput model.User
		updateError error

		hasJwtError bool
	}{
		{
			name:     "standart input with ID=110",
			jwtToken: testingJwtToken,

			updtaeInput: testingModelUser,
			updateError: nil,

			hasJwtError: false,
		},
		{
			name:     "input with repository update error",
			jwtToken: testingJwtToken,

			updtaeInput: testingModelUser,
			updateError: errSome,

			hasJwtError: false,
		},
		{
			name:     "input with jwt assertion error",
			jwtToken: jwt.New(jwt.SigningMethodHS256),

			updtaeInput: testingModelUser,
			updateError: nil,

			hasJwtError: true,
		},
	}
	for _, test := range testTable {
		c := e.NewContext(req, rec)
		c.Set("user", test.jwtToken)
		ctx := c.Request().Context()

		updateCall := user.EXPECT().Update(ctx, test.updtaeInput).Return(test.updateError).Maybe()
		err := handler.UpdateUser(c)

		if test.hasJwtError {
			assert.Error(t, err, test.name)
			return
		}

		user.AssertExpectations(t)
		user.AssertCalled(t, "Update", ctx, test.updtaeInput)
		updateCall.Unset()
	}
}

func TestDeleteUser(t *testing.T) { // aditional testing required; reason: repository delete do not return error when deleting raw's tha do not exist
	user := mocks.NewUserInterface(t)
	password := mocks.NewPasswordInterface(t)
	handler := NewHandler(user, password, nil)

	method := http.MethodDelete
	target := "/users/auth/delete"

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, nil)

	testTable := []struct {
		name     string
		jwtToken *jwt.Token
		id       int

		hasJwtError         bool
		deleteUserError     error
		deletePasswordError error
	}{
		{
			name:        "jwt token assertion error",
			jwtToken:    jwt.New(jwt.SigningMethodHS256),
			hasJwtError: true,
		},
		{
			name:     "standart input without error & ID=110",
			jwtToken: testingJwtToken,
			id:       110,

			hasJwtError:         false,
			deleteUserError:     nil,
			deletePasswordError: nil,
		},
		{
			name:     "input with deleteUser error & ID=110",
			jwtToken: testingJwtToken,
			id:       110,

			hasJwtError:     false,
			deleteUserError: errSome,
		},
		{
			name:     "input with deletePassword error & ID=110",
			jwtToken: testingJwtToken,
			id:       110,

			hasJwtError:         false,
			deleteUserError:     nil,
			deletePasswordError: errSome,
		},
	}

	for _, test := range testTable {
		c := e.NewContext(req, rec)
		c.Set("user", test.jwtToken) // inconsequential
		ctx := c.Request().Context()

		usrDeleteCall := user.EXPECT().Delete(ctx, test.id).Return(test.deleteUserError).Maybe()
		pwdDeleteCall := password.EXPECT().DeletePassword(ctx, test.id).Return(test.deletePasswordError).Maybe()
		err := handler.DeleteUser(c)

		if test.hasJwtError {
			assert.Error(t, err, test.name)
			return
		}

		assert.Nil(t, err, test.name)

		user.AssertExpectations(t)
		user.AssertCalled(t, "Delete", ctx, test.id)

		if test.deleteUserError == nil {
			password.AssertExpectations(t)
			password.AssertCalled(t, "DeletePassword", ctx, test.id)
		}

		usrDeleteCall.Unset()
		pwdDeleteCall.Unset()
	}

}
