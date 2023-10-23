package handlers

import (
	"github.com/labstack/echo"
	"github.com/mmfshirokan/GoProject1/model"
	"github.com/mmfshirokan/GoProject1/service"

	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	password *service.Password
	user     *service.User
	err      error
}

func NewHandler(usr *service.User, usrpw *service.Password) *Handler {
	return &Handler{
		user:     usr,
		password: usrpw,
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
	return c.String(http.StatusOK, "Usser id: "+strconv.FormatInt(int64(usr.Id), 10)+"\nUser name: "+usr.Name+"\nUser male:"+strconv.FormatBool(usr.Male)+"\n")
}

/*func (hand *Handler) createUser(c echo.Context) error {
	var usr model.User
	hand.err = c.Bind(&usr)
	if hand.err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	hand.err = hand.user.Create(usr.Id, usr.Name, usr.Male)
	return hand.err
}*/

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

func (hand *Handler) Register(c echo.Context) error {
	var usr model.User
	hand.err = c.Bind(&usr)
	if hand.err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	hand.err = hand.user.Create(usr.Id, usr.Name, usr.Male)
	if hand.err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	password := c.FormValue("password")

	hand.err = hand.password.Store(usr.Id, password)
	return hand.err
}

func (hand *Handler) Login(password string, login string, c echo.Context) (bool, error) {
	var usr model.User
	hand.err = c.Bind(&usr)
	if hand.err != nil {
		return false, c.String(http.StatusBadRequest, "bad request")
	}
	password = c.FormValue("password")
	var correct bool

	correct, hand.err = hand.password.Compare(usr.Id, password)
	return correct, hand.err
}
