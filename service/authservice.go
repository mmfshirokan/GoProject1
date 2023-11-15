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
	"github.com/mmfshirokan/GoProject1/handlers/request"
	"github.com/mmfshirokan/GoProject1/model"
	"github.com/mmfshirokan/GoProject1/repository"
)

type Token struct {
	repo repository.AuthRepositoryInterface
}

func NewToken(rep repository.AuthRepositoryInterface) *Token {
	return &Token{
		repo: rep,
	}
}

func (tok *Token) CreateAuthToken(usr model.User) string {
	claims := &request.UserRequest{
		Id:   usr.Id,
		Name: usr.Name,
		Male: usr.Male,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
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
	id := uuid.New()

	hashed_id, err := tok.conductHashing(id)
	if err != nil {
		return fmt.Errorf("hashing: %w", err)
	}

	return tok.repo.Create(ctx, &model.RefreshToken{
		UserID:     userID,
		ID:         id,
		Hash:       hashed_id,
		Expiration: time.Now().Add(time.Hour * 12),
	})
}

func (tok *Token) ValidateRfTokenTrougId(recived_hash string, id uuid.UUID) (bool, error) {
	//add time validation
	expected_hash, err := tok.conductHashing(id)
	if err != nil {
		return false, fmt.Errorf("hashing: %w", err)
	}

	res := expected_hash == recived_hash
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

func (tok *Token) GetByUserID(ctx context.Context, userId int) ([]*model.RefreshToken, error) {
	return tok.repo.GetByUserID(ctx, userId)
}
