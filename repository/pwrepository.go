package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PwRepositoryInterface interface {
	Store(uint, string) error
	Compare(uint, string) (bool, error)
}

func NewPasswordRepository(conf config.Config) PwRepositoryInterface {
	if conf.Database == "mongodb" {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:6543"))
		defer client.Disconnect(context.Background())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect client: %v\n", err)
		}

		collection := client.Database("users").Collection("passwords")

		return &repositoryMongo{
			client:     client,
			collection: collection,
		}
	}

	dbpool, err := pgxpool.New(context.Background(), " postgres://echopguser:pgpw4echo@localhost:8080/echopwdb?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS passwords (id INT PRIMARY KEY, password CHARACTER VARYING(30) NOT NULL)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table in PostgresDB: %v\n", err)
	}

	return &repositoryPostgres{
		dbpool: dbpool,
	}
}

func (rep *repositoryMongo) Store(id uint, pw string) error {
	_, err := rep.collection.InsertOne(context.Background(), bson.D{
		{Key: "_id", Value: id},
		{Key: "password", Value: pw},
	})
	return err
}

func (rep *repositoryMongo) Compare(id uint, pw string) (bool, error) {
	var dbpw string
	err := rep.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&dbpw)
	if dbpw == pw {
		return true, err
	}
	return false, err
}

func (rep *repositoryPostgres) Store(id uint, pw string) error {
	_, err := rep.dbpool.Exec(context.Background(), "INSERT INTO passwords VALUES ($1, $2)", id, pw)
	return err
}

func (rep *repositoryPostgres) Compare(id uint, pw string) (bool, error) {
	var dbpw string
	err := rep.dbpool.QueryRow(context.Background(), "SELECT password FROM passwords WHERE id = $1", id).Scan(&dbpw)
	if dbpw == pw {
		return true, err
	}
	return false, err
}
