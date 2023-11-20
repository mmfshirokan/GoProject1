package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/model"
	"github.com/mmfshirokan/GoProject1/service"

	"net/http"
	"strconv"
)

type Handler struct {
	password *service.Password
	user     *service.User
	token    *service.Token
}

func NewHandler(usr *service.User, usrpw *service.Password, tok *service.Token) *Handler {
	return &Handler{
		user:     usr,
		password: usrpw,
		token:    tok,
	}
}

func (hand *Handler) GetUser(c echo.Context) error {
	c.Request().Context()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.UserRequest)
	return c.String(http.StatusOK, "Usser id: "+strconv.FormatInt(int64(claims.Id), 10)+"\nUser name: "+claims.Name+"\nUser male:"+strconv.FormatBool(claims.Male)+"\n")
}

func (hand *Handler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.UserRequest)
	err := hand.user.Update(ctx, claims.Id, claims.Name, claims.Male)
	return err
}

func (hand *Handler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.UserRequest)
	if err := hand.user.Delete(ctx, claims.Id); err != nil {
		return err
	}
	err := hand.password.DeletePassword(ctx, claims.Id)
	return err
}
