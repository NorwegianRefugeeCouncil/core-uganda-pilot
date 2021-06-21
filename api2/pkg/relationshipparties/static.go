package relationshipparties

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(ctx context.Context, partiesStore *PartiesStore) error {
	if _, err := partiesStore.store.Collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{
				"attributes."+attributes.FirstNameAttribute.ID, "text",
			},
			{
				"attributes."+attributes.LastNameAttribute.ID, "text",
			},
		},
	}); err != nil {
		return err
	}
	return nil
}