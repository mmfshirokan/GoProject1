package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/config"
	"github.com/mmfshirokan/GoProject1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"fmt"
	"os"
)

type RepositoryInterface interface {
	GetTroughID(context.Context, int) (string, bool, error) //TODO стоит заменить int, strin, bool на model
	Update(context.Context, int, string, bool) error
	Create(context.Context, int, string, bool) error
	Delete(context.Context, int) error
}

type repositoryMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type repositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NewRepository(conf config.Config) RepositoryInterface {
	if conf.Database == "mongodb" {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:6543"))
		//defer client.Disconnect(context.Background())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect client: %v\n", err)
		}

		collection := client.Database("users").Collection("entity")

		return &repositoryMongo{
			client:     client,
			collection: collection,
		}
	}

	dbpool, err := pgxpool.New(context.Background(), "postgres://echopguser:pgpw4echo@localhost:5432/echodb?sslmode=disable") //postgres://echopguser:pgadminpwd4echo@localhost:5432/echodb?sslmode=disable// os.Getenv("DATABASE_URL")
	if err != nil {
		dbpool.Close()
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil
	}

	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS entity (id INT PRIMARY KEY, name TEXT NOT NULL, male BOOLEAN NOT NULL)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table in PostgresDB: %v\n", err)
	}

	return &repositoryPostgres{
		dbpool: dbpool,
	}
}

func (rep *repositoryMongo) GetTroughID(ctx context.Context, id int) (string, bool, error) {
	var usr model.User
	err := rep.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&usr)
	return usr.Name, usr.Male, err
}

func (rep *repositoryMongo) Create(ctx context.Context, id int, name string, male bool) error {
	_, err := rep.collection.InsertOne(ctx, bson.D{
		{Key: "_id", Value: id},
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	return err
}

func (rep *repositoryMongo) Update(ctx context.Context, id int, name string, male bool) error {
	_, err := rep.collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	return err
}

func (rep *repositoryMongo) Delete(ctx context.Context, id int) error {
	_, err := rep.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

func (rep *repositoryPostgres) GetTroughID(ctx context.Context, id int) (string, bool, error) {
	usr := model.User{}
	err := rep.dbpool.QueryRow(ctx, "SELECT name, male FROM entity WHERE id = $1", id).Scan(&usr.Name, &usr.Male)
	return usr.Name, usr.Male, err
}

func (rep *repositoryPostgres) Create(ctx context.Context, id int, name string, male bool) error {
	_, err := rep.dbpool.Exec(ctx, "INSERT INTO entity VALUES ($1, $2, $3)", id, name, male)
	return err
}

func (rep *repositoryPostgres) Update(ctx context.Context, id int, name string, male bool) error {
	_, err := rep.dbpool.Exec(ctx, "UPDATE entity SET name = $1, male = $2 WHERE id = $3", name, male, id)
	return err
}

func (rep *repositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := rep.dbpool.Exec(ctx, "DELETE FROM entity WHERE id = $1", id)
	return err
}
