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

// GetUser godoc
//
// @Summary get user
// @Description Gets stored user data
// @Tags User Handlers
// @Produce json
// @Param token header string true "JWT token"
// @Success 200 {object} model.User
// @Router /users/auth/get [get]
func (handling *Handler) GetUser(c echo.Context) error {
	logInit()

	c.Request().Context()

	token, okey := c.Get("user").(*jwt.Token)
	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims, okey := token.Claims.(*model.UserRequest)
	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(http.StatusBadRequest)
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
		log.Error(fmt.Errorf("handlers.GetUser; c.String: %w", err))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

// UpdateUser godoc
//
// @Summary update user
// @Description Updates user data
// @Tags User Handlers
// @Accept json
// @Param token header string true "JWT token"
// @Success 1
// @Router /users/auth/update [put]
func (handling *Handler) UpdateUser(c echo.Context) error {
	logInit()

	ctx := c.Request().Context()
	token, okey := c.Get("user").(*jwt.Token)

	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims, okey := token.Claims.(*model.UserRequest)
	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := handling.user.Update(ctx, model.User{
		ID:   claims.ID,
		Name: claims.Name,
		Male: claims.Male,
	}); err != nil {
		log.Error(fmt.Errorf("handling.user.Update: %w", err))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

// DeleteUser godoc
//
// @Summary delete user
// @Description Deletes user data
// @Tags User Handlers
// @Param token header string true "JWT token"
// @Success 1
// @Router /users/auth/delete [delete]
func (handling *Handler) DeleteUser(c echo.Context) error {
	logInit()

	ctx := c.Request().Context()
	token, okey := c.Get("user").(*jwt.Token)

	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims, okey := token.Claims.(*model.UserRequest)
	if !okey {
		log.Error(fmt.Errorf("jwt agregation failed"))

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := handling.user.Delete(ctx, claims.ID); err != nil {
		log.Error(fmt.Errorf("handling.user.Delete: %w", err))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := handling.password.DeletePassword(ctx, claims.ID); err != nil {
		log.Error(fmt.Errorf("handling.password.DeletePassword: %w", err))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func logInit() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.TextFormatter{})
}
