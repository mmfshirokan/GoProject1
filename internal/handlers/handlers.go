package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/service"
	log "github.com/sirupsen/logrus"
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
	logInit()

	c.Request().Context()

	token, okey := c.Get("user").(*jwt.Token)
	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(400)
	}

	claims, okey := token.Claims.(*model.UserRequest)
	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(400)
	}

	err := c.String(http.StatusOK, fmt.Sprint(
		"Usser id: ",
		strconv.FormatInt(int64(claims.ID), 10),
		"\nUser name: ",
		claims.Name,
		"\nUser male:",
		strconv.FormatBool(claims.Male),
		"\n",
	))
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(400, err.Error())
	}

	return nil
}

func (handling *Handler) UpdateUser(c echo.Context) error {
	logInit()

	ctx := c.Request().Context()
	token, okey := c.Get("user").(*jwt.Token)

	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(400)
	}

	claims, okey := token.Claims.(*model.UserRequest)
	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(400)
	}

	if err := handling.user.Update(ctx, claims.ID, claims.Name, claims.Male); err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(400, err.Error())
	}

	return nil
}

func (handling *Handler) DeleteUser(c echo.Context) error {
	logInit()

	ctx := c.Request().Context()
	token, okey := c.Get("user").(*jwt.Token)

	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(400)
	}

	claims, okey := token.Claims.(*model.UserRequest)
	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(400)
	}

	if err := handling.user.Delete(ctx, claims.ID); err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(400, err.Error())
	}

	if err := handling.password.DeletePassword(ctx, claims.ID); err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(400, err.Error())
	}

	return nil
}

func logInit() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.TextFormatter{})
}
