package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/model"
	"github.com/mmfshirokan/GoProject1/service"
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

func (handling *Handler) GetUser(c echo.Context) error {
	c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil
	}

	claims, ok := token.Claims.(*model.UserRequest)
	if !ok {
		return nil
	}

	return fmt.Errorf("%w", c.String(http.StatusOK, fmt.Sprint(
		"Usser id: ",
		strconv.FormatInt(int64(claims.ID), 10),
		"\nUser name: ",
		claims.Name,
		"\nUser male:",
		strconv.FormatBool(claims.Male),
		"\n",
	)))
}

func (handling *Handler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)

	if !ok {
		return nil
	}

	claims, ok := token.Claims.(*model.UserRequest)
	if !ok {
		return nil
	}

	err := handling.user.Update(ctx, claims.ID, claims.Name, claims.Male)

	return fmt.Errorf("user.Update: %w", err)
}

func (handling *Handler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)

	if !ok {
		return nil
	}

	claims, ok := token.Claims.(*model.UserRequest)
	if !ok {
		return nil
	}

	if err := handling.user.Delete(ctx, claims.ID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	err := handling.password.DeletePassword(ctx, claims.ID)

	return fmt.Errorf("%w", err)
}
