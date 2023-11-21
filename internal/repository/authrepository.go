package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/internal/model"
)

type AuthRepositoryInterface interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByUserID(ctx context.Context, userID int) ([]*model.RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type authRepositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NewAuthRpository() AuthRepositoryInterface {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, "postgres://echopguser:pgpw4echo@localhost:5432/echodb?sslmode=disable")
	if err != nil {
		dbpool.Close()
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)

		return nil
	}

	_, err = dbpool.Exec(context.Background(), fmt.Sprint(
		"CREATE TABLE IF NOT EXISTS rf_tokens ",
		"(user_id INT NOT NULL, id UUID PRIMARY KEY, hash TEXT NOT NULL, expire TIME NOT NULL)",
	))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table in PostgresDB: %v\n", err)
	}

	return &authRepositoryPostgres{
		dbpool: dbpool,
	}
}

func (rep *authRepositoryPostgres) Create(ctx context.Context, token *model.RefreshToken) error {
	const maxNumberOfTokens = 5

	reshreshTokens, _ := rep.GetByUserID(ctx, token.UserID)

	if len(reshreshTokens) > maxNumberOfTokens {
		if err := rep.Delete(ctx, reshreshTokens[0].ID); err != nil {
			return err
		}
	}

	err := rep.create(ctx, token)

	return fmt.Errorf("rep.creat: %w", err)
}

func (rep *authRepositoryPostgres) GetByUserID(ctx context.Context, userID int) ([]*model.RefreshToken, error) {
	rows, err := rep.dbpool.Query(ctx, fmt.Sprint(
		"SELECT id, user_id, hash, expire ",
		"FROM rf_tokens WHERE user_id = $1 order by expire desc",
	), userID)
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
	_, err := rep.dbpool.Exec(ctx, "DELETE FROM rf_tokens WHERE id = $1", id)

	return fmt.Errorf("dpool.Exec: %w", err)
}

func (rep *authRepositoryPostgres) create(ctx context.Context, token *model.RefreshToken) error {
	_, err := rep.dbpool.Exec(ctx, fmt.Sprint(
		"INSERT INTO rf_tokens ",
		"VALUES ($1, $2, $3, $4)",
	), token.UserID, token.ID, token.Hash, token.Expiration)

	return fmt.Errorf("dpool.Exe: %w", err)
}
