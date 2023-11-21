package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/internal/model"
)

func (handling *Handler) SignUp(echoContext echo.Context) error {
	var usr model.User

	err := echoContext.Bind(&usr)
	if err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	ctx := echoContext.Request().Context()

	err = handling.user.Create(ctx, usr.ID, usr.Name, usr.Male)
	if err != nil {
		return fmt.Errorf("create method in handAuth: %w", err)
	}

	err = handling.password.Store(ctx, usr.ID, usr.Password)
	if err != nil {
		return fmt.Errorf("create method in handAuth: %w", err)
	}

	return nil
}

func (handling *Handler) SignIn(echoContext echo.Context) error {
	var usr model.User
	if err := echoContext.Bind(&usr); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	ctx := echoContext.Request().Context()

	validPassword, err := handling.password.Compare(ctx, usr.ID, usr.Password)
	if err != nil {
		return fmt.Errorf("password.Compare: %w", err)
	}

	if !validPassword {
		return nil
	}

	authToken := handling.token.CreateAuthToken(usr.ID, usr.Name, usr.Male)
	err = handling.token.CreateRfToken(ctx, usr.ID)

	if err != nil {
		return fmt.Errorf("token.createRfToken: %w", err)
	}

	var refeshTokens []*model.RefreshToken

	refeshTokens, err = handling.token.GetByUserID(ctx, usr.ID)
	if err != nil {
		return fmt.Errorf("token.GetTokenTroughId: %w", err)
	}

	return fmt.Errorf(
		"json: %w", echoContext.JSON(http.StatusOK, echo.Map{
			"token":   authToken,
			"refresh": refeshTokens[0],
		}),
	)
}

func (handling *Handler) Refresh(echoContext echo.Context) error {
	var refreshToken model.RefreshToken
	if err := echoContext.Bind(&refreshToken); err != nil {
		return fmt.Errorf("binding model.rf_token: %w", err)
	}

	ctx := echoContext.Request().Context()

	valid, err := handling.token.ValidateRfTokenTrougID(refreshToken.Hash, refreshToken.ID)
	if err != nil {
		return fmt.Errorf("error ocured while validating rf_token: %w", err)
	}

	if !valid {
		if err = handling.token.Delete(ctx, refreshToken.ID); err != nil {
			return fmt.Errorf("token.Delete: %w", err)
		}

		return nil
	}

	err = handling.token.Delete(ctx, refreshToken.ID)
	if err != nil {
		return fmt.Errorf("token.Delete: %w", err)
	}

	err = handling.token.CreateRfToken(ctx, refreshToken.UserID)
	if err != nil {
		return fmt.Errorf("token.CreateRefreshToken: %w", err)
	}

	var rfTokens []*model.RefreshToken

	rfTokens, err = handling.token.GetByUserID(ctx, refreshToken.UserID)
	if err != nil {
		return fmt.Errorf("token.GetByUserID: %w", err)
	}

	name, male, err := handling.user.GetTroughID(ctx, refreshToken.UserID)
	if err != nil {
		return fmt.Errorf("usr.GetTrouhId: %w", err)
	}

	return fmt.Errorf(
		"json: %w", echoContext.JSON(http.StatusOK, echo.Map{
			"token":   handling.token.CreateAuthToken(refreshToken.UserID, name, male),
			"refresh": rfTokens[0],
		}),
	)
}
