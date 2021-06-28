package testing

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func ClearDatabase(ctx context.Context, databaseName string, mongoClient *mongo.Client) error {
	return mongoClient.Database(databaseName).Drop(ctx)
}
