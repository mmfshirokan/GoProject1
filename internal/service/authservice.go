package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
)

type Token struct {
	repo repository.AuthRepositoryInterface
}

func NewToken(rep repository.AuthRepositoryInterface) *Token {
	return &Token{
		repo: rep,
	}
}

func (tok *Token) CreateAuthToken(id int, name string, male bool) string {
	const authTokenLifeTime = 6
	claims := &model.UserRequest{
		ID:   id,
		Name: name,
		Male: male,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * authTokenLifeTime)),
		},
	}

	authTok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	result, err := authTok.SignedString([]byte("secret"))
	if err != nil {
		return ""
	}

	return result
}

func (tok *Token) CreateRfToken(ctx context.Context, userID int) error {
	const refreshTokenLifeTime = 12

	id := uuid.New()

	hashedID, err := tok.conductHashing(id)
	if err != nil {
		return fmt.Errorf("authService.conductHasing in authService.CreateRfToken: %w", err)
	}

	return tok.repo.Create(ctx, &model.RefreshToken{
		UserID:     userID,
		ID:         id,
		Hash:       hashedID,
		Expiration: time.Now().Add(time.Hour * refreshTokenLifeTime),
	})
}

func (tok *Token) ValidateRfTokenTrougID(receivedHash string, id uuid.UUID) (bool, error) {
	expectedHash, err := tok.conductHashing(id)
	if err != nil {
		return false, fmt.Errorf("authService.conductHasing in authService.ValidateRfTokenTrougID: %w", err)
	}

	res := (expectedHash == receivedHash)

	return res, nil
}

func (tok *Token) conductHashing(id uuid.UUID) (string, error) {
	h := hmac.New(sha256.New, []byte("secret"))

	marsheled, err := json.Marshal(id)
	if err != nil {
		return "", fmt.Errorf("json.Marsha: %w", err)
	}

	str := base64.URLEncoding.EncodeToString(marsheled)

	_, err = h.Write([]byte(str))
	if err != nil {
		return "", fmt.Errorf("json.Marsha: %w", err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil)), nil
}

func (tok *Token) Delete(ctx context.Context, id uuid.UUID) error {
	return tok.repo.Delete(ctx, id)
}

func (tok *Token) GetByUserID(ctx context.Context, userID int) ([]*model.RefreshToken, error) {
	return tok.repo.GetByUserID(ctx, userID)
}
