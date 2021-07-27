package attachments

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttachmentStore struct {
	getCollection utils.MongoCollectionFn
}

func NewAttachmentStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*AttachmentStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &AttachmentStore{
		getCollection: utils.GetCollectionFn(database, "attachments", mongoClientFn),
	}

	collection, err := store.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return store, nil
}
