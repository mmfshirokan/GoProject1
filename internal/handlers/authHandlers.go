package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/service"
	log "github.com/sirupsen/logrus"
)

// Signup godoc
//
// @Summary signup
// @Description Adding new user to the database
// @Tags User Authentication
// @Accept json
// @Param usr body model.User true "User Data"
// @Success 200
// @Router /users/signup [post]
func (handling *Handler) SignUp(echoContext echo.Context) error {
	logInit()

	var usr model.User

	err := echoContext.Bind(&usr)
	if err != nil {
		log.Error(fmt.Errorf("binding erorr at handlers.SignUp: %w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	err = val.Struct(&usr)

	if err != nil {
		log.Error(fmt.Errorf("model.User struct validation error at handlers.SignUP: %w", err))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := echoContext.Request().Context()

	err = handling.user.Create(ctx, usr)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = handling.password.Store(ctx, usr)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

// Signin godoc
//
// @Summary signin
// @Description Comparing user input password with database password and giving him access to other handlers
// @Tags User Authentication
// @Accept json
// @Produce json
// @Param usr body model.User true "Only Password and ID required"
// @Success 200 {object} map[string]interface{}
// @Router /users/signin [put]
func (handling *Handler) SignIn(echoContext echo.Context) error {
	logInit()

	var usr model.User
	if err := echoContext.Bind(&usr); err != nil {
		log.Error(fmt.Errorf("binding erorr: %w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(&usr); err != nil {
		log.Error(fmt.Errorf("validation error at model.User: %w", err))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := echoContext.Request().Context()

	validPassword, err := handling.password.Compare(ctx, usr)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !validPassword {
		err = errors.New("invalid password for sign in method")
		log.Error(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	authToken := service.CreateAuthToken(usr.ID, usr.Name, usr.Male)
	err = handling.token.CreateRfToken(ctx, usr.ID)

	if err != nil {
		log.Error(fmt.Errorf("token.createRfToken: %w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var refeshTokens []*model.RefreshToken

	refeshTokens, err = handling.token.GetByUserID(ctx, usr.ID)
	if err != nil {
		log.Error(fmt.Errorf("token.GetTokenTroughId: %w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = echoContext.JSON(http.StatusOK, echo.Map{
		"token":   authToken,
		"refresh": refeshTokens[0],
	})
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

// Refresh godoc
//
// @Summary refresh
// @Description Refreshes token paire
// @Tags User Authentication
// @Accept json
// @Produce json
// @Param refreshToken body model.RefreshToken true "Refresh Token Data"
// @Success 200 {object} map[string]interface{}
// @Router /users/refresh [put]
func (handling *Handler) Refresh(echoContext echo.Context) error {
	logInit()

	var refreshToken model.RefreshToken
	if err := echoContext.Bind(&refreshToken); err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(&refreshToken); err != nil {
		log.Error(fmt.Errorf("validation error at model.RfToken: %w", err))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := echoContext.Request().Context()

	valid, err := service.ValidateRfTokenTroughID(refreshToken.Hash, refreshToken.ID)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !valid {
		if err = handling.token.Delete(ctx, refreshToken.ID); err != nil {
			log.Error(fmt.Errorf("%w", err))

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return nil
	}

	err = handling.token.Delete(ctx, refreshToken.ID)
	if err != nil {
		log.Error(fmt.Errorf("token.Delete: %w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = handling.token.CreateRfToken(ctx, refreshToken.UserID)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var rfTokens []*model.RefreshToken

	rfTokens, err = handling.token.GetByUserID(ctx, refreshToken.UserID)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	usr, err := handling.user.GetTroughID(ctx, refreshToken.UserID)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = echoContext.JSON(http.StatusOK, echo.Map{
		"token":   service.CreateAuthToken(refreshToken.UserID, usr.Name, usr.Male),
		"refresh": rfTokens[0],
	})
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
