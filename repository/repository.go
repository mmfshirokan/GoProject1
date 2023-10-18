package repository

import (
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

type Repository struct {
	client     *mongo.Client
	collection *mongo.Collection
	err        error
}

func NewRepository() *Repository {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:5432"))
	defer client.Disconnect(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect client: %v\n", err)
		os.Exit(1)
	}

	collection := client.Database("echodb").Collection("users")

	return &Repository{
		client:     client,
		collection: collection,
		err:        err,
	}
}

func (rep *Repository) GetTroughID(id int) (string, bool, error) {
	var usr model.User
	rep.err = rep.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&usr)
	return usr.Name, usr.Male, rep.err
}

func (rep *Repository) Create(id int, name string, male bool) error {
	_, rep.err = rep.collection.InsertOne(context.Background(), bson.D{
		{Key: "_id", Value: id},
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	return rep.err
}

func (rep *Repository) Update(id int, name string, male bool) error {
	_, rep.err = rep.collection.ReplaceOne(context.Background(), bson.D{{Key: "_id", Value: id}}, bson.D{
		{Key: "name", Value: name},
		{Key: "male", Value: male},
	})
	return rep.err
}

func (rep *Repository) Delete(id int) error {
	_, rep.err = rep.collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	return rep.err
}
