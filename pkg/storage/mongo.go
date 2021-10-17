package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

type Factory interface {
	New() (*mongo.Client, error)
}

type MongoFn func() (*mongo.Client, error)

type factory struct {
	mongoFn MongoFn
}

func (f factory) New() (*mongo.Client, error) {
	return f.mongoFn()
}

func NewFactory(fn MongoFn) Factory {
	return factory{mongoFn: fn}
}

func ClearCollections(ctx context.Context, mongoCli *mongo.Client, databaseName string, collectionNames ...string) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, collectionName := range collectionNames {
		c := collectionName
		g.Go(func() error {
			_, err := mongoCli.Database(databaseName).Collection(c).DeleteMany(ctx, bson.M{})
			return err
		})
	}
	return g.Wait()
}
