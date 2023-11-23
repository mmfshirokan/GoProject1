package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/internal/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PwRepositoryInterface interface {
	Store(ctx context.Context, id int, pw string) error
	Compare(ctx context.Context, id int, pw string) (bool, error)
	DeletePassword(ctx context.Context, id int) error
}

func NewPasswordRepository(conf config.Config) PwRepositoryInterface {
	ctx := context.Background()

	if conf.Database == "mongodb" {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoURL))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect client: %v\n", err)
		}

		collection := client.Database("users").Collection("passwords")

		return &repositoryMongo{
			client:     client,
			collection: collection,
		}
	}

	dbpool, err := pgxpool.New(ctx, conf.PostgresURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	}

	return &repositoryPostgres{
		dbpool: dbpool,
	}
}

func (rep *repositoryMongo) Store(ctx context.Context, id int, pw string) error {
	logInit()

	log.WithFields(log.Fields{
		"id":       id,
		"password": pw,
	}).Debug("method: repository.Store")

	_, err := rep.collection.InsertOne(ctx, bson.D{
		{Key: "_id", Value: id},
		{Key: "password", Value: pw},
	})
	if err != nil {
		return fmt.Errorf("InsertOne in repository.Store%w", err)
	}

	return nil
}

func (rep *repositoryMongo) Compare(ctx context.Context, id int, pw string) (bool, error) {
	logInit()

	log.WithFields(log.Fields{
		"id":       id,
		"password": pw,
	}).Debug("method: repository.Compare")

	var dbpw string

	err := rep.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&dbpw)
	if err != nil {
		return false, fmt.Errorf("findOne.Decode in repository.Compare: %w", err)
	}

	if dbpw == pw {
		return true, nil
	}

	return false, nil
}

func (rep *repositoryMongo) DeletePassword(ctx context.Context, id int) error {
	logInit()

	log.WithField("id", id).Debug("method: repository.DeletePassword")

	_, err := rep.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return fmt.Errorf("deleteOne in repository.DeletePassword: %w", err)
	}

	return nil
}

func (rep *repositoryPostgres) Store(ctx context.Context, id int, pw string) error {
	logInit()

	log.WithFields(log.Fields{
		"id":       id,
		"password": pw,
	}).Error("method: repository.Store")

	_, err := rep.dbpool.Exec(ctx, "INSERT INTO apps.passwords VALUES ($1, $2)", id, pw)
	if err != nil {
		return fmt.Errorf("exec in repository.Store%w", err)
	}

	return nil
}

func (rep *repositoryPostgres) Compare(ctx context.Context, id int, pw string) (bool, error) {
	logInit()

	log.WithFields(log.Fields{
		"id":       id,
		"password": pw,
	}).Debug("method: repository.Compare")

	var dbpw string

	err := rep.dbpool.QueryRow(ctx, "SELECT password FROM apps.passwords WHERE id = $1", id).Scan(&dbpw)
	if err != nil {
		return false, fmt.Errorf("queryRow.Scan in repository.Compare: %w", err)
	}

	if dbpw == pw {
		return true, nil
	}

	return false, nil
}

func (rep *repositoryPostgres) DeletePassword(ctx context.Context, id int) error {
	logInit()

	log.WithField("id", id).Error("method: repository.DeletePassword")

	_, err := rep.dbpool.Exec(ctx, "DELETE FROM apps.passwords WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec in repository.DeletePassword%w", err)
	}

	return nil
}

func logInit() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}
