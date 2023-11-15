package repository

import (
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/model"

	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	//"github.com/mmfshirokan/GoProject1/config"

	"context"
	"fmt"
)

type AuthRepositoryInterface interface {
	Create(context.Context, *model.RefreshToken) error
	GetByUserID(context.Context, int) ([]*model.RefreshToken, error)
	Delete(context.Context, uuid.UUID) error
}

type authRepositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NewAuthRpository() AuthRepositoryInterface {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		//fmt.Errorf("pool conection: %w", err)
		return nil
	}
	return &authRepositoryPostgres{
		dbpool: dbpool,
	}
}

func (rep *authRepositoryPostgres) Create(ctx context.Context, token *model.RefreshToken) error {
	reshreshTokens, _ := rep.GetByUserID(ctx, token.UserID)

	if len(reshreshTokens) > 5 {
		if err := rep.Delete(ctx, reshreshTokens[0].ID); err != nil {
			return err
		}
	}
	err := rep.create(ctx, token)
	return err
}

func (rep *authRepositoryPostgres) GetByUserID(ctx context.Context, userId int) ([]*model.RefreshToken, error) {
	rows, err := rep.dbpool.Query(ctx, "SELECT id, user_id, hash, expire  FROM rf_tokens WHERE user_id = $1 order by expire desc", userId)
	if err != nil {
		return make([]*model.RefreshToken, 0), fmt.Errorf("query %w", err)
	}
	defer rows.Close()
	retsult := make([]*model.RefreshToken, 0)

	for rows.Next() {
		item := &model.RefreshToken{}
		err := rows.Scan(&item.ID, &item.UserID, &item.Hash, &item.Expiration)
		if err != nil {

			return make([]*model.RefreshToken, 0), fmt.Errorf("scan %w", err)
		}
		retsult = append(retsult, item)
	}

	return retsult, nil
}

func (rep *authRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := rep.dbpool.Exec(ctx, "DELETE FROM rf_tokens WHERE user_id = $1", id)
	return err
}

func (rep *authRepositoryPostgres) create(ctx context.Context, token *model.RefreshToken) error {
	_, err := rep.dbpool.Exec(ctx, "INSERT INTO rf_tokens VALUES ($1, $2, $3, $4)", token.UserID, token.ID, token.Hash, token.Expiration)
	return err
}
