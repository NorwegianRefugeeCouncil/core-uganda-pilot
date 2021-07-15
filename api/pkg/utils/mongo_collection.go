package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoClientFn = func(ctx context.Context) (*mongo.Client, error)
type MongoCollectionFn = func(ctx context.Context) (*mongo.Collection, error)

func GetCollectionFn(database, collection string, clientFn MongoClientFn) MongoCollectionFn {
	return func(ctx context.Context) (*mongo.Collection, error) {
		mongoClient, err := clientFn(ctx)
		if err != nil {
			return nil, err
		}
		return mongoClient.Database(database).Collection(collection), nil
	}
}
