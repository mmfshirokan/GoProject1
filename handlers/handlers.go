package handlers

import (
	"fmt"

	"net/http"

	"github.com/mmfshirokan/GoProject1/model"

	"github.com/mmfshirokan/GoProject1/service"

	"github.com/labstack/echo"
)

type Handler struct {
	serv *service.Service
	err  error
}

func NewHandler() *Handler {
	hand := Handler{
		serv: service.NewService(),
	}
	return &hand
}

func (hand *Handler) GetUser(c echo.Context) error {
	usr := model.User{
		Id:   c.Param("id"),
		Name: c.FormValue("name"),
		Male: c.FormValue("male"),
	}

	usr.Name, usr.Male, hand.err = hand.serv.GetUserTroughID(usr.Id)
	if hand.err != nil {
		fmt.Println("Error ocured: ", hand.err)
	}
	return c.String(http.StatusOK, "Usser id: "+usr.Id+"\nUser name: "+usr.Name+"\nUser male:"+usr.Male+"\n")
}

func (hand *Handler) SaveUser(c echo.Context) error {
	usr := model.User{
		Id:   c.Param("id"),
		Name: c.FormValue("name"),
		Male: c.FormValue("male"),
	}
	hand.err = hand.serv.SaveUser(usr.Id, usr.Name, usr.Male)
	return hand.err
}

func (hand *Handler) UpdateUser(c echo.Context) error {
	usr := model.User{
		Id:   c.Param("id"),
		Name: c.FormValue("name"),
		Male: c.FormValue("male"),
	}
	hand.err = hand.serv.UpdateUser(usr.Id, usr.Name, usr.Male)
	return hand.err
}

func (hand *Handler) DeleteUser(c echo.Context) error {
	usr := model.User{
		Id: c.Param("id"),
	}
	hand.err = hand.serv.DeleteUser(usr.Id)
	return hand.err
}
