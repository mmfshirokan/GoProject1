package handlers

import (
	"net/http"

	"repository"

	"github.com/labstack/echo"
)

type User struct {
	id   string
	name string
	male string
}

type Handler struct {
}

func NewUserHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetUser(c echo.Context) error {
	usr := User{
		id:   c.Param("id"),
		name: c.FormValue("name"),
		male: c.FormValue("male"),
	}
	usr.male, usr.male, _ = h.serv.GetUserTroughID(usr.id)
	return c.String(http.StatusOK, "User id: "+usr.id+"\nUser name: "+usr.name+"\nUser male: "+usr.male+"\n")
}

func SaveUser(c echo.Context) error {
	usr := User{
		id:   c.Param("id"),
		name: c.FormValue("name"),
		male: c.FormValue("male"),
	}
	err := repository.SaveUser(usr.id, usr.name, usr.male)
	return err
}

func UpdateUser(c echo.Context) error {
	usr := User{
		id:   c.Param("id"),
		name: c.FormValue("name"),
		male: c.FormValue("male"),
	}
	err := repository.UpdateUser(usr.id, usr.name, usr.male)
	return err
}

func DeleteUser(c echo.Context) error {
	usr := User{
		id: c.Param("id"),
	}
	err := repository.DeleteUser(usr.id)
	return err
}
