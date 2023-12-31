package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/internal/model"
	log "github.com/sirupsen/logrus"
)

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

	err = handling.user.Create(ctx, usr.ID, usr.Name, usr.Male)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = handling.password.Store(ctx, usr.ID, usr.Password)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

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

	validPassword, err := handling.password.Compare(ctx, usr.ID, usr.Password)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !validPassword {
		log.Error(fmt.Errorf("invalid password for sign in method"))

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	authToken := handling.token.CreateAuthToken(usr.ID, usr.Name, usr.Male)
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

	valid, err := handling.token.ValidateRfTokenTrougID(refreshToken.Hash, refreshToken.ID)
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

	name, male, err := handling.user.GetTroughID(ctx, refreshToken.UserID)
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = echoContext.JSON(http.StatusOK, echo.Map{
		"token":   handling.token.CreateAuthToken(refreshToken.UserID, name, male),
		"refresh": rfTokens[0],
	})
	if err != nil {
		log.Error(fmt.Errorf("%w", err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
