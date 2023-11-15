package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

	rfTokens, err := handling.token.GetByUserID(ctx, int(usr.Id))
	if err != nil {
		return fmt.Errorf("GetByUserId: %w", err)
	}

	//Возможно стоит упростить эту часть?
	var (
		cook                  *http.Cookie
		ifCookieExistsItIsNil error
	)
	for _, val := range rfTokens {
		cook, ifCookieExistsItIsNil = c.Cookie(val.ID.String())

		if ifCookieExistsItIsNil == nil {
			break
		}
	}

	var validToken bool

	tempUUIDstring, tempUUIDerr := uuid.Parse(cook.Name)
	if tempUUIDerr != nil {
		return fmt.Errorf("WrongUUID: %w", tempUUIDerr)
	}

	validToken, err = handling.token.ValidateRfTokenTrougId(cook.Value, tempUUIDstring)
	if tempUUIDerr != nil {
		return fmt.Errorf("WrongVAlidation: %w", err)
	}

	if ifCookieExistsItIsNil != nil || !validToken {
		validPassword, err := handling.password.Compare(ctx, usr.Id, usr.Password)
		if err != nil {
			return fmt.Errorf("password.Compare: %w", err)
		}

		if !validPassword {
			return fmt.Errorf("invalid password")
		}
	}

	autht := handling.token.CreateAuthToken(usr)
	err = handling.token.CreateRfToken(ctx, int(usr.Id))

	if err != nil {
		return fmt.Errorf("token.createRfToken: %w", err)
	}

	var rft []*model.RefreshToken
	rft, err = handling.token.GetByUserID(ctx, int(usr.Id))
	if err != nil {
		return fmt.Errorf("token.GetTokenTroughId: %w", err)
	}

	c.SetCookie(handling.setRfCookie(rft[0]))
	return c.JSON(http.StatusOK, echo.Map{
		"token": autht,
	})
}

func (handling *Handler) setRfCookie(value *model.RefreshToken) *http.Cookie {
	return &http.Cookie{
		Name:    value.ID.String(),
		Value:   value.Hash,
		Expires: value.Expiration,
		Path:    "/users/auth",
	}
}
