package repository

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockery --dir . --all --output ../service/mocks --with-expecter

type RepositoryInterface interface {
	GetTroughID(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, usr model.User) error
	Create(ctx context.Context, usr model.User) error
	Delete(ctx context.Context, id int) error
}

type repositoryMongo struct {
	collection *mongo.Collection
}

type repositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NewPostgresRepository(dbpool *pgxpool.Pool) RepositoryInterface {
	return &repositoryPostgres{
		dbpool: dbpool,
	}
}

func NewMongoRepository(collection *mongo.Collection) RepositoryInterface {
	return &repositoryMongo{
		collection: collection,
	}
}

func (rep *repositoryMongo) GetTroughID(ctx context.Context, id int) (*model.User, error) { //nolint:gocritic // it is unconvinient to name results because of decode
	usr := &model.User{}

	err := rep.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&usr)
	if err != nil {
		return nil, fmt.Errorf("findOne.Decode in repository.GetTroughID: %w", err)
	}

	return usr, nil
}

func (rep *repositoryMongo) Create(ctx context.Context, usr model.User) error {
	_, err := rep.collection.InsertOne(ctx, bson.D{
		{Key: "_id", Value: usr.ID},
		{Key: "name", Value: usr.Name},
		{Key: "male", Value: usr.Male},
	})
	if err != nil {
		return fmt.Errorf("insertOne in repository.Create: %w", err)
	}

	return nil
}

func (rep *repositoryMongo) Update(ctx context.Context, usr model.User) error {
	_, err := rep.collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: usr.ID}}, bson.D{
		{Key: "name", Value: usr.Name},
		{Key: "male", Value: usr.Male},
	})
	if err != nil {
		return fmt.Errorf("replaceOne in repository.Update: %w", err)
	}

	return nil
}

func (rep *repositoryMongo) Delete(ctx context.Context, id int) error {
	_, err := rep.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return fmt.Errorf("deleteOne in repository.Delete: %w", err)
	}

	return nil
}

func (rep *repositoryPostgres) GetTroughID(ctx context.Context, id int) (*model.User, error) { //nolint:gocritic // it is unconvinient to name results because of decode
	usr := &model.User{ID: id}
	err := rep.dbpool.QueryRow(ctx, "SELECT name, male FROM apps.entity WHERE id = $1", id).Scan(&usr.Name, &usr.Male)
	if err != nil {
		return nil, fmt.Errorf("queryRow in repository.GetTroughID: %w", err)
	}

	return usr, nil
}

func (rep *repositoryPostgres) Create(ctx context.Context, usr model.User) error {
	if err := validate(usr); err != nil {
		return err
	}

	_, err := rep.dbpool.Exec(ctx, "INSERT INTO apps.entity VALUES ($1, $2, $3)", usr.ID, usr.Name, usr.Male)
	if err != nil {
		return fmt.Errorf("exec in repository.Create: %w", err)
	}

	return nil
}

func (rep *repositoryPostgres) Update(ctx context.Context, usr model.User) error {
	if err := validate(usr); err != nil {
		return err
	}

	_, err := rep.dbpool.Exec(ctx, "UPDATE apps.entity SET name = $1, male = $2 WHERE id = $3", usr.Name, usr.Male, usr.ID)
	if err != nil {
		return fmt.Errorf("exec in repository.Update: %w", err)
	}

	return nil
}

func (rep *repositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := rep.dbpool.Exec(ctx, "DELETE FROM apps.entity WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec in repository.Delete: %w", err)
	}

	return nil
}

func validate(usr model.User) error {
	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(&usr); err != nil {
		return err
	}

	return nil
}
