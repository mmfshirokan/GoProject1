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
	Store(context.Context, int, string) error
	Compare(context.Context, int, string) (bool, error)
	DeletePassword(context.Context, int) error
}

func NewPasswordRepository(conf config.Config) PwRepositoryInterface {
	ctx := context.Background() //should it also be some special context?
	if conf.Database == "mongodb" {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:6543"))
		//defer client.Disconnect(context.Background())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect client: %v\n", err)
		}

		collection := client.Database("users").Collection("passwords")

		return &repositoryMongo{
			client:     client,
			collection: collection,
		}
	}

	dbpool, err := pgxpool.New(ctx, "postgres://echopguser:pgpw4echo@localhost:5432/echodb?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	}

	_, err = dbpool.Exec(ctx, "CREATE TABLE IF NOT EXISTS passwords (id INT PRIMARY KEY, password TEXT NOT NULL)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table in PostgresDB: %v\n", err)
	}

	return &repositoryPostgres{
		dbpool: dbpool,
	}
}

func (rep *repositoryMongo) Store(ctx context.Context, id int, pw string) error {
	_, err := rep.collection.InsertOne(ctx, bson.D{
		{Key: "_id", Value: id},
		{Key: "password", Value: pw},
	})
	return err
}

func (rep *repositoryMongo) Compare(ctx context.Context, id int, pw string) (bool, error) {
	var dbpw string
	err := rep.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&dbpw)
	if dbpw == pw {
		return true, err
	}
	return false, err
}

func (rep *repositoryMongo) DeletePassword(ctx context.Context, id int) error {
	_, err := rep.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

func (rep *repositoryPostgres) Store(ctx context.Context, id int, pw string) error {
	_, err := rep.dbpool.Exec(ctx, "INSERT INTO passwords VALUES ($1, $2)", id, pw)
	return err
}

func (rep *repositoryPostgres) Compare(ctx context.Context, id int, pw string) (bool, error) {
	var dbpw string
	err := rep.dbpool.QueryRow(ctx, "SELECT password FROM passwords WHERE id = $1", id).Scan(&dbpw)
	if dbpw == pw {
		return true, err
	}
	return false, err
}

func (rep *repositoryPostgres) DeletePassword(ctx context.Context, id int) error {
	_, err := rep.dbpool.Exec(ctx, "DELETE FROM passwords WHERE id = $1", id)
	return err
}
