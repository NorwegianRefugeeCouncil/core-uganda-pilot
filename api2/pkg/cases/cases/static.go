package cases

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DomesticAbuse = Case{
	ID:          "763bcc10-444a-4853-b3d6-13a1224a77ec",
	CaseTypeID:  "2ab2aba2-c18f-4235-9ccd-52161defca5b",
	PartyID:     "ab7a1620-f34e-4811-8534-853167ed7944",
	Description: "Domestic abuse",
	Done:        false,
}
var MonthlyAllowance = Case{
	ID:          "47499762-c189-4a74-9156-7969f899073b",
	CaseTypeID:  "4b37e5d0-56e7-48b3-8227-bed8ce72019a",
	PartyID:     "40b30fb0-c392-4798-9400-bda3e5837867",
	Description: "Monthly allowance",
	Done:        false,
}

var ChildCare = Case{
	ID:          "8fb5f755-85eb-4d91-97a9-fdf86c01df25",
	CaseTypeID:  "73f47b43-eaa3-4ece-af91-0a72ff4c742e",
	PartyID:     "40b30fb0-c392-4798-9400-bda3e5837867",
	Description: "Monthly stipend for Bo Diddley's child",
	Done:        true,
}

var Eviction = Case{
	ID:          "1189c79a-e1df-4af8-8f0b-19619970d410",
	CaseTypeID:  "9a5ee26f-8df3-447c-a4b0-ed7f36710d95",
	PartyID:     "0bde06f0-5416-4514-9c5a-794a2cc2f1b7",
	Description: "Support in eviction case",
	Done:        false,
}

func Init(ctx context.Context, store *Store) error {
	if _, err := store.collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return err
	}

	for _, kase := range []Case{
		DomesticAbuse,
		MonthlyAllowance,
		ChildCare,
		Eviction,
	} {
		if err := store.Create(ctx, &kase); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := store.Update(ctx, &kase); err != nil {
				return err
			}
		}
	}
	return nil
}
