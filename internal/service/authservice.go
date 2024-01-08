package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
)

type TokenInterface interface {
	CreateAuthToken(id int, name string, male bool) string
	CreateRfToken(ctx context.Context, userID int) error
	ValidateRfTokenTroughID(receivedHash string, id uuid.UUID) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserID(ctx context.Context, userID int) ([]*model.RefreshToken, error)
}

type Token struct {
	repo      repository.AuthRepositoryInterface
	redis     repository.RedisRepositoryInterface[[]*model.RefreshToken]
	sourceMap repository.MapRepositoryInterface[[]*model.RefreshToken]
}

func NewToken(
	rep repository.AuthRepositoryInterface,
	redis repository.RedisRepositoryInterface[[]*model.RefreshToken],
	sourceMap repository.MapRepositoryInterface[[]*model.RefreshToken],
) TokenInterface {
	return &Token{
		repo:      rep,
		redis:     redis,
		sourceMap: sourceMap,
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

	hashedID, err := conductHashing(id)
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

func (tok *Token) ValidateRfTokenTroughID(receivedHash string, id uuid.UUID) (bool, error) {
	expectedHash, err := conductHashing(id)
	if err != nil {
		return false, fmt.Errorf("authService.conductHasing in authService.ValidateRfTokenTrougID: %w", err)
	}

	res := (expectedHash == receivedHash)

	return res, nil
}

func (tok *Token) Delete(ctx context.Context, id uuid.UUID) error {
	return tok.repo.Delete(ctx, id)
}

func (tok *Token) GetByUserID(ctx context.Context, userID int) ([]*model.RefreshToken, error) {
	key := "token:" + strconv.FormatInt(int64(userID), 10)
	mod, err := tok.sourceMap.Get(ctx, key)

	if err != nil {
		mod, err = tok.repo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		err = tok.redis.Add(ctx, key, mod)
		if err != nil {
			return mod, fmt.Errorf("redis XAdd error at repository redis: %w", err)
		}
	}

	return mod, nil
}

func conductHashing(id uuid.UUID) (string, error) {
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
