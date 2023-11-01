package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mmfshirokan/GoProject1/handlers/request"

	"net/http"
	"strconv"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type Token struct {
}

func (tok *Token) GenerateAuthToken(id uint, name string, male bool) (string, error) {
	claims := &request.UserRequest{
		Id:   id,
		Name: name,
		Male: male,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, err
}

func (tok *Token) GenerateRefreshToken(id uint) *http.Cookie {
	strId := strconv.FormatUint(uint64(id), 10)
	key := []byte("secret")
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strId))

	cookie := &http.Cookie{
		Name:     strconv.FormatUint(uint64(id), 10),
		Value:    base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil)),
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 72),
	}
	return cookie
}

func (tok *Token) ValidateRefreshToken(id uint, rfToken string) bool {
	cook := tok.GenerateRefreshToken(id)
	if cook.Value == rfToken {
		return true
	}
	return false
}
