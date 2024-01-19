package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmfshirokan/GoProject1/internal/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PwRepositoryInterface interface {
	Store(ctx context.Context, usr model.User) error
	Compare(ctx context.Context, usr model.User) (bool, error)
	DeletePassword(ctx context.Context, id int) error
	BulkStore(ctx context.Context, pwd [][]interface{}) error
}

func NewPostgresPasswordRepository(dbpool *pgxpool.Pool) PwRepositoryInterface {
	return &repositoryPostgres{
		dbpool: dbpool,
	}
}

func NewMongoPasswordRepository(collection *mongo.Collection) PwRepositoryInterface {
	return &repositoryMongo{
		collection: collection,
	}
}

func (rep *repositoryMongo) Store(ctx context.Context, usr model.User) error {
	logInit()

	log.WithFields(log.Fields{
		"id":       usr.ID,
		"password": usr.Password,
	}).Debug("method: repository.Store")

	_, err := rep.collection.InsertOne(ctx, bson.D{
		{Key: "_id", Value: usr.ID},
		{Key: "password", Value: usr.Password},
	})
	if err != nil {
		return fmt.Errorf("InsertOne in repository.Store%w", err)
	}

	return nil
}

func (rep *repositoryMongo) Compare(ctx context.Context, usr model.User) (bool, error) {
	logInit()

	log.WithFields(log.Fields{
		"id":       usr.ID,
		"password": usr.Password,
	}).Debug("method: repository.Compare")

	var dbpw string

	err := rep.collection.FindOne(ctx, bson.D{{Key: "_id", Value: usr.ID}}).Decode(&dbpw)
	if err != nil {
		return false, fmt.Errorf("findOne.Decode in repository.Compare: %w", err)
	}

	if dbpw == usr.Password {
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

func (rep *repositoryMongo) BulkStore(ctx context.Context, pwd [][]interface{}) error {
	err := errors.New("Not implemented exeption")
	return err
}

func (rep *repositoryPostgres) Store(ctx context.Context, usr model.User) error {
	logInit()

	log.WithFields(log.Fields{
		"id":       usr.ID,
		"password": usr.Password,
	}).Debug("method: repository.Store")

	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(&usr); err != nil {
		log.Error("validation error ocured:", err)
		return err
	}

	_, err := rep.dbpool.Exec(ctx, "INSERT INTO apps.passwords VALUES ($1, $2)", usr.ID, usr.Password)
	if err != nil {
		log.Error("exec in repository.Store: %w", err)
		return err
	}

	return nil
}

func (rep *repositoryPostgres) Compare(ctx context.Context, usr model.User) (bool, error) {
	logInit()

	log.WithFields(log.Fields{
		"id":       usr.ID,
		"password": usr.Password,
	}).Debug("method: repository.Compare")

	var dbpw string

	err := rep.dbpool.QueryRow(ctx, "SELECT password FROM apps.passwords WHERE id = $1", usr.ID).Scan(&dbpw)
	if err != nil {
		return false, fmt.Errorf("queryRow.Scan in repository.Compare: %w", err)
	}

	if dbpw == usr.Password {
		return true, nil
	}

	return false, nil
}

func (rep *repositoryPostgres) DeletePassword(ctx context.Context, id int) error {
	logInit()

	log.WithField("id", id).Debug("method: repository.DeletePassword")

	_, err := rep.dbpool.Exec(ctx, "DELETE FROM apps.passwords WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec in repository.DeletePassword%w", err)
	}

	return nil
}

// Only for Kafka:
func (rep *repositoryPostgres) BulkStore(ctx context.Context, pwd [][]interface{}) error {
	_, err := rep.dbpool.CopyFrom(ctx, pgx.Identifier{"apps", "passwords"}, []string{"id", "password"}, pgx.CopyFromRows(pwd))
	log.Info("Bulk Store")
	if err != nil {
		return err
	}
	return nil
}

func logInit() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}
