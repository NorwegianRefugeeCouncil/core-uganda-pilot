package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoClientFn = func(ctx context.Context) (*mongo.Client, error)
type MongoCollectionFn = func(ctx context.Context) (collection *mongo.Collection, close func(), err error)

func GetCollectionFn(database, collection string, clientFn MongoClientFn) MongoCollectionFn {
	return func(ctx context.Context) (*mongo.Collection, func(), error) {
		mongoClient, err := clientFn(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("failed to get mongo client")
			return nil, nil, err
		}
		return mongoClient.Database(database).Collection(collection), func() {
			mongoClient.Disconnect(ctx)
		}, nil
	}
}
