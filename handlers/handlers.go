package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/handlers/request"
	"github.com/mmfshirokan/GoProject1/model"
	"github.com/mmfshirokan/GoProject1/service"

	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	password *service.Password
	user     *service.User
}

func NewHandler(usr *service.User, usrpw *service.Password) *Handler {
	return &Handler{
		user:     usr,
		password: usrpw,
	}
}

func (hand *Handler) GetUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*request.UserRequest)
	return c.String(http.StatusOK, "Usser id: "+strconv.FormatInt(int64(claims.Id), 10)+"\nUser name: "+claims.Name+"\nUser male:"+strconv.FormatBool(claims.Male)+"\n")
}

func (hand *Handler) UpdateUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*request.UserRequest)
	var user model.User
	c.Bind(&user)
	err := hand.user.Update(claims.Id, user.Name, user.Male)
	return err
}

func (hand *Handler) DeleteUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*request.UserRequest)
	err := hand.user.Delete(claims.Id)
	return err
}

func (hand *Handler) Register(c echo.Context) error {
	var usr model.User
	err := c.Bind(&usr)
	if err != nil {
		return err
	}
	err = hand.user.Create(usr.Id, usr.Name, usr.Male)
	if err != nil {
		return err
	}

	err = hand.password.Store(usr.Id, usr.Password)
	return err
}

func (hand *Handler) Login(c echo.Context) error {
	var usr model.User
	if err := c.Bind(&usr); err != nil {
		return err
	}
	correct, err := hand.password.Compare(usr.Id, usr.Password)
	if !correct || err != nil {
		return echo.ErrUnauthorized
	}

	usr.Name, usr.Male, err = hand.user.GetTroughID(usr.Id)
	if err != nil {
		return err
	}

	claims := &request.UserRequest{
		Id:   usr.Id,
		Name: usr.Name,
		Male: usr.Male,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
