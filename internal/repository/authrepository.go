package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/model"
)

type AuthRepositoryInterface interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByUserID(ctx context.Context, userID int) ([]*model.RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type authRepositoryPostgres struct {
	dbpool   *pgxpool.Pool
	redisrep *repositoryRedis[[]*model.RefreshToken]
}

func NewAuthRpository(conf config.Config) AuthRepositoryInterface {
	if conf.Database == "mongodb" { //nolint:goconst //unnecessary const
		return nil
	}

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, conf.PostgresURI)
	if err != nil {
		dbpool.Close()
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)

		return nil
	}

	return &authRepositoryPostgres{
		dbpool:   dbpool,
		redisrep: NewRfTokenRedisRepository(conf),
	}
}

func (rep *authRepositoryPostgres) Create(ctx context.Context, token *model.RefreshToken) error {
	const maxNumberOfTokens = 5

	reshreshTokens, _ := rep.GetByUserID(ctx, token.UserID)

	if len(reshreshTokens) > maxNumberOfTokens {
		if err := rep.Delete(ctx, reshreshTokens[0].ID); err != nil {
			return fmt.Errorf("authRepository.Delte in authRepository.Create: %w", err)
		}
	}

	if err := rep.create(ctx, token); err != nil {
		return fmt.Errorf("authRepository.create in authRepository.Create")
	}

	return nil
}

func (rep *authRepositoryPostgres) GetByUserID(ctx context.Context, userID int) ([]*model.RefreshToken, error) {
	mod, err := rep.redisrep.Get(ctx, strconv.FormatInt(int64(userID), 10))

	if err != nil {
		rows, err := rep.dbpool.Query(ctx, fmt.Sprint(
			"SELECT id, user_id, hash, expire ",
			"FROM apps.rf_tokens WHERE user_id = $1 order by expire desc",
		), userID)

		if err != nil {
			return make([]*model.RefreshToken, 0), fmt.Errorf("query in GetByUserID: %w", err)
		}
		defer rows.Close()

		retsult := make([]*model.RefreshToken, 0)

		for rows.Next() {
			item := &model.RefreshToken{}

			err := rows.Scan(&item.ID, &item.UserID, &item.Hash, &item.Expiration)
			if err != nil {
				return make([]*model.RefreshToken, 0), fmt.Errorf("rows.Scan in authRepository.GetByUSerID: %w", err)
			}

			retsult = append(retsult, item)
		}

		err = rep.redisrep.Set(ctx, strconv.FormatInt(int64(userID), 10), retsult)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}

		return retsult, nil
	}

	return mod, nil
}

func (rep *authRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := rep.dbpool.Exec(ctx, "DELETE FROM apps.rf_tokens WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec in authRepository.Delete")
	}

	return nil
}

func (rep *authRepositoryPostgres) create(ctx context.Context, token *model.RefreshToken) error {
	_, err := rep.dbpool.Exec(ctx, fmt.Sprint(
		"INSERT INTO apps.rf_tokens ",
		"VALUES ($1, $2, $3, $4)",
	), token.UserID, token.ID, token.Hash, token.Expiration)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
