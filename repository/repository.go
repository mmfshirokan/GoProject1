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
	GetTroughID(uint) (string, bool, error)
	Update(uint, string, bool) error
	Create(uint, string, bool) error
	Delete(uint) error
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
		defer client.Disconnect(context.Background())

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
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS entity (id INT PRIMARY KEY, name CHARACTER VARYING(30) NOT NULL, male BOOLEAN NOT NULL)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table in PostgresDB: %v\n", err)
	}

	return &repositoryPostgres{
		dbpool: dbpool,
	}
}

func (rep *repositoryMongo) GetTroughID(id uint) (string, bool, error) {
	var usr model.User
	err := rep.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&usr)
	return usr.Name, usr.Male, err
}

func (rep *repositoryMongo) Create(id uint, name string, male bool) error {
	_, err := rep.collection.InsertOne(context.Background(), bson.D{
		{Key: "_id", Value: id},
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	return err
}

func (rep *repositoryMongo) Update(id uint, name string, male bool) error {
	_, err := rep.collection.ReplaceOne(context.Background(), bson.D{{Key: "_id", Value: id}}, bson.D{
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	return err
}

func (rep *repositoryMongo) Delete(id uint) error {
	_, err := rep.collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	return err
}

func (rep *repositoryPostgres) GetTroughID(id uint) (string, bool, error) {
	usr := model.User{}
	err := rep.dbpool.QueryRow(context.Background(), "SELECT name, male FROM entity WHERE id = $1", id).Scan(&usr.Name, &usr.Male)
	return usr.Name, usr.Male, err
}

func (rep *repositoryPostgres) Create(id uint, name string, male bool) error {
	_, err := rep.dbpool.Exec(context.Background(), "INSERT INTO entity VALUES ($1, $2, $3)", id, name, male)
	return err
}

func (rep *repositoryPostgres) Update(id uint, name string, male bool) error {
	_, err := rep.dbpool.Exec(context.Background(), "UPDATE entity SET name = $1, male = $2 WHERE id = $3", name, male, id)
	return err
}

func (rep *repositoryPostgres) Delete(id uint) error {
	_, err := rep.dbpool.Exec(context.Background(), "DELETE FROM entity WHERE id = $1", id)
	return err
}
