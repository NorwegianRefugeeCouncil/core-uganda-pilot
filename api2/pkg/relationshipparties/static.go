package relationshipparties

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(ctx context.Context, partiesStore *PartiesStore) error {
	if _, err := partiesStore.store.Collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"$**": "text"},
	}); err != nil {
		return err
	}
	return nil
}