package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/internal/config"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Interface interface {
	GetTroughID(ctx context.Context, id int) (string, bool, error)
	Update(ctx context.Context, id int, name string, male bool) error
	Create(ctx context.Context, id int, name string, male bool) error
	Delete(ctx context.Context, id int) error
}

type repositoryMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type repositoryPostgres struct {
	dbpool   *pgxpool.Pool
	redisrep *repositoryRedis[model.User]
}

func NewRepository(conf config.Config) Interface {
	if conf.Database == "mongodb" {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(conf.MongoURI))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect client: %v\n", err)
		}

		collection := client.Database("users").Collection("entity")

		return &repositoryMongo{
			client:     client,
			collection: collection,
		}
	}

	dbpool, err := pgxpool.New(context.Background(), conf.PostgresURI)
	if err != nil {
		dbpool.Close()
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)

		return nil
	}

	return &repositoryPostgres{
		dbpool:   dbpool,
		redisrep: NewUserRedisRepository(conf),
	}
}

func (rep *repositoryMongo) GetTroughID(ctx context.Context, id int) (string, bool, error) { //nolint:gocritic // it is unconvinient to name results because of decode
	var usr model.User

	err := rep.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&usr)
	if err != nil {
		return "", false, fmt.Errorf("findOne.Decode in repository.GetTroughID: %w", err)
	}

	return usr.Name, usr.Male, nil
}

func (rep *repositoryMongo) Create(ctx context.Context, id int, name string, male bool) error {
	_, err := rep.collection.InsertOne(ctx, bson.D{
		{Key: "_id", Value: id},
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	if err != nil {
		return fmt.Errorf("insertOne in repository.Create: %w", err)
	}

	return nil
}

func (rep *repositoryMongo) Update(ctx context.Context, id int, name string, male bool) error {
	_, err := rep.collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{
		{Key: "name", Value: name},
		{Key: "male", Value: male},
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

func (rep *repositoryPostgres) GetTroughID(ctx context.Context, id int) (string, bool, error) { //nolint:gocritic // it is unconvinient to name results because of decode
	usr := model.User{}
	usr, err := rep.redisrep.Get(ctx, strconv.FormatInt(int64(id), 10))

	if err != nil {
		err := rep.dbpool.QueryRow(ctx, "SELECT name, male FROM apps.entity WHERE id = $1", id).Scan(&usr.Name, &usr.Male)
		if err != nil {
			return "", false, fmt.Errorf("queryRow in repository.GetTroughID: %w", err)
		}

		err = rep.redisrep.Set(ctx, strconv.FormatInt(int64(id), 10), model.User{
			ID:   id,
			Name: usr.Name,
			Male: usr.Male,
		})
		if err != nil {
			return "", false, fmt.Errorf("%w", err)
		}
	}

	return usr.Name, usr.Male, nil
}

func (rep *repositoryPostgres) Create(ctx context.Context, id int, name string, male bool) error {
	_, err := rep.dbpool.Exec(ctx, "INSERT INTO apps.entity VALUES ($1, $2, $3)", id, name, male)
	if err != nil {
		return fmt.Errorf("exec in repository.Create: %w", err)
	}

	return nil
}

func (rep *repositoryPostgres) Update(ctx context.Context, id int, name string, male bool) error {
	_, err := rep.dbpool.Exec(ctx, "UPDATE apps.entity SET name = $1, male = $2 WHERE id = $3", name, male, id)
	if err != nil {
		return fmt.Errorf("exec in repository.Update: %w", err)
	}

	err = rep.redisrep.Remove(ctx, strconv.FormatInt(int64(id), 10))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (rep *repositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := rep.dbpool.Exec(ctx, "DELETE FROM apps.entity WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec in repository.Delete: %w", err)
	}

	err = rep.redisrep.Remove(ctx, strconv.FormatInt(int64(id), 10))
	if err != nil {
		fmt.Fprintf(os.Stderr, "redis hash wasn't there %v", id)
	}

	return nil
}
