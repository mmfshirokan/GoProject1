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
	GetTroughID(int) (string, bool, error)
	Update(int, string, bool) error
	Create(int, string, bool) error
	Delete(int) error
}

type RepositoryMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
	err        error
}

type RepositoryPostgres struct {
	dbpool *pgxpool.Pool
	err    error
}

func NewRepository(conf config.Config) RepositoryInterface {
	if conf.Database == "mongodb" {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:6543"))
		defer client.Disconnect(context.Background())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect client: %v\n", err)
		}

		collection := client.Database("users").Collection("entity")

		return &RepositoryMongo{
			client:     client,
			collection: collection,
			err:        err,
		}
	}

	dbpool, err := pgxpool.New(context.Background(), " postgres://echopguser:pgadminpwd4echo@localhost:8080/echodb?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS entity (id INT PRIMARY KEY, name CHARACTER VARYING(30) NOT NULL, male BOOLEAN NOT NULL)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table in PostgresDB: %v\n", err)
	}

	return &RepositoryPostgres{
		dbpool: dbpool,
		err:    err,
	}
}

func (rep *RepositoryMongo) GetTroughID(id int) (string, bool, error) {
	var usr model.User
	rep.err = rep.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&usr)
	return usr.Name, usr.Male, rep.err
}

func (rep *RepositoryMongo) Create(id int, name string, male bool) error {
	_, rep.err = rep.collection.InsertOne(context.Background(), bson.D{
		{Key: "_id", Value: id},
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	return rep.err
}

func (rep *RepositoryMongo) Update(id int, name string, male bool) error {
	_, rep.err = rep.collection.ReplaceOne(context.Background(), bson.D{{Key: "_id", Value: id}}, bson.D{
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	return rep.err
}

func (rep *RepositoryMongo) Delete(id int) error {
	_, rep.err = rep.collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	return rep.err
}

func (rep *RepositoryPostgres) GetTroughID(id int) (string, bool, error) {
	usr := model.User{}
	rep.err = rep.dbpool.QueryRow(context.Background(), "SELECT name, male FROM entity WHERE id = $1", id).Scan(&usr.Name, &usr.Male)
	return usr.Name, usr.Male, rep.err
}

func (rep *RepositoryPostgres) Create(id int, name string, male bool) error {
	_, rep.err = rep.dbpool.Exec(context.Background(), "INSERT INTO entity VALUES ($1, $2, $3)", id, name, male)
	return rep.err
}

func (rep *RepositoryPostgres) Update(id int, name string, male bool) error {
	_, rep.err = rep.dbpool.Exec(context.Background(), "UPDATE entity SET name = $1, male = $2 WHERE id = $3", name, male, id)
	return rep.err
}

func (rep *RepositoryPostgres) Delete(id int) error {
	_, rep.err = rep.dbpool.Exec(context.Background(), "DELETE FROM entity WHERE id = $1", id)
	return rep.err
}
