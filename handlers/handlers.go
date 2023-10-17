package handlers

import (
	"fmt"

	"os"

	"net/http"

	"github.com/mmfshirokan/GoProject1/model"

	"github.com/mmfshirokan/GoProject1/service"

	"github.com/labstack/echo"
)

type Handler struct {
	user *service.User
	err  error
}

func NewHandler(usr *service.User) *Handler {
	return &Handler{
		user: usr,
	}
}

func (hand *Handler) GetUser(c echo.Context) error {
	var usr model.User
	hand.err = c.Bind(&usr)
	if hand.err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	usr.Name, usr.Male, hand.err = hand.user.GetTroughID(usr.Id)
	if hand.err != nil {
		fmt.Fprintf(os.Stderr, "Error ocured while operating with db: %v\n", hand.err)
		return hand.err
	}
	return c.String(http.StatusOK, "Usser id: "+usr.Id+"\nUser name: "+usr.Name+"\nUser male:"+usr.Male+"\n")
}

func (hand *Handler) SaveUser(c echo.Context) error {
	var usr model.User
	hand.err = c.Bind(&usr)
	if hand.err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	hand.err = hand.user.Create(usr.Id, usr.Name, usr.Male)
	return hand.err
}

func (hand *Handler) UpdateUser(c echo.Context) error {
	var usr model.User
	hand.err = c.Bind(&usr)
	if hand.err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	hand.err = hand.user.Update(usr.Id, usr.Name, usr.Male)
	return hand.err
}

func (hand *Handler) DeleteUser(c echo.Context) error {
	var usr model.User
	hand.err = c.Bind(&usr)
	if hand.err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	hand.err = hand.user.Delete(usr.Id)
	return hand.err
}
