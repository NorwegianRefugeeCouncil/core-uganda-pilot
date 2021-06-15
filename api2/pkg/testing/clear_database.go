package testing

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func ClearDatabase(ctx context.Context, mongoClient *mongo.Client) error {
	return mongoClient.Database("core").Drop(ctx)
}
