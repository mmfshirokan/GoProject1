package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/handlers/request"
	"github.com/mmfshirokan/GoProject1/service"

	"net/http"
	"strconv"
)

type Handler struct {
	password *service.Password
	user     *service.User
	token    *service.Token
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
	err := hand.user.Update(claims.Id, claims.Name, claims.Male)
	return err
}

func (hand *Handler) DeleteUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*request.UserRequest)
	if err := hand.user.Delete(claims.Id); err != nil {
		return err
	}
	err := hand.password.DeletePassword(claims.Id)
	return err
}
