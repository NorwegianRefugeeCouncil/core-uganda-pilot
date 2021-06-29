package testhelpers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(ctx context.Context) (*mongo.Client, error) {

	mongoClient, err := mongo.NewClient(
		options.Client().ApplyURI("mongodb://localhost:27017"),
		options.Client().SetAuth(options.Credential{
			Username: "root",
			Password: "example",
		}),
	)
	if err != nil {
		return nil, err
	}
	if err := mongoClient.Connect(ctx); err != nil {
		return nil, err
	}
	return mongoClient, nil
}
