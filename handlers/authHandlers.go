package handlers

import (
	"fmt"
	"net/http"

	//"github.com/google/uuid"
	//"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"github.com/labstack/echo/v4"
	"github.com/mmfshirokan/GoProject1/model"
)

func (handling *Handler) SignUp(c echo.Context) error {
	var usr model.User
	err := c.Bind(&usr)
	if err != nil {
		return fmt.Errorf("bind: %w", err)
	}
	ctx := c.Request().Context()

	err = handling.user.Create(ctx, usr.Id, usr.Name, usr.Male)
	if err != nil {
		return fmt.Errorf("create method in handAuth: %w", err)
	}

	err = handling.password.Store(ctx, usr.Id, usr.Password)
	if err != nil {
		return fmt.Errorf("create method in handAuth: %w", err)
	}

	return nil
}

func (handling *Handler) SignIn(c echo.Context) error {
	var usr model.User
	if err := c.Bind(&usr); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	ctx := c.Request().Context()

	validPassword, err := handling.password.Compare(ctx, usr.Id, usr.Password)
	if err != nil {
		return fmt.Errorf("password.Compare: %w", err)
	}

	if !validPassword {
		return fmt.Errorf("invalid password")
	}

	authToken := handling.token.CreateAuthToken(usr.Id, usr.Name, usr.Male)
	err = handling.token.CreateRfToken(ctx, usr.Id)

	if err != nil {
		return fmt.Errorf("token.createRfToken: %w", err)
	}

	var refeshTokens []*model.RefreshToken
	refeshTokens, err = handling.token.GetByUserID(ctx, usr.Id)
	if err != nil {
		return fmt.Errorf("token.GetTokenTroughId: %w", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":   authToken,
		"refresh": refeshTokens[0],
	})
}

func (handling *Handler) Refresh(c echo.Context) error {
	var refresh_token model.RefreshToken
	if err := c.Bind(&refresh_token); err != nil {
		return fmt.Errorf("binding model.rf_token: %w", err)
	}
	ctx := c.Request().Context()
	valid, err := handling.token.ValidateRfTokenTrougId(refresh_token.Hash, refresh_token.ID)
	if err != nil {
		return fmt.Errorf("error ocured while validating rf_token: %w", err)
	}
	if !valid {
		handling.token.Delete(ctx, refresh_token.ID)
		return fmt.Errorf("rf_token not valid")
	}
	err = handling.token.Delete(ctx, refresh_token.ID)
	if err != nil {
		return err
	}
	err = handling.token.CreateRfToken(ctx, refresh_token.UserID)
	if err != nil {
		return err
	}
	var rf_tokens []*model.RefreshToken
	rf_tokens, err = handling.token.GetByUserID(ctx, refresh_token.UserID)
	if err != nil {
		return err
	}

	name, male, err := handling.user.GetTroughID(ctx, refresh_token.UserID)
	if err != nil {
		return fmt.Errorf("usr.GetTrouhId: %w", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":   handling.token.CreateAuthToken(refresh_token.UserID, name, male),
		"refresh": rf_tokens[0],
	})
}
