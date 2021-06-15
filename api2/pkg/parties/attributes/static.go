package attributes

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var FirstNameAttribute = Attribute{
	ID:   "8514da51-aad5-4fb4-a797-8bcc0c969b27",
	Name: "firstName",
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "First Name",
			ShortFormulation: "First Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypes: []string{
		partytypes.IndividualPartyType.ID,
	},
}

var LastNameAttribute = Attribute{
	ID:   "21079bbc-e04b-4fe8-897f-644d73af0d9e",
	Name: "lastName",
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Last Name",
			ShortFormulation: "Last Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypes: []string{
		partytypes.IndividualPartyType.ID,
	},
}

var BirthDateAttribute = Attribute{
	ID:   "87fe07d7-e6a7-4428-8086-3842b69f3665",
	Name: "birthDate",
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Birth Date",
			ShortFormulation: "Birth Date",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypes: []string{
		partytypes.IndividualPartyType.ID,
	},
}

func Init(ctx context.Context, store *Store) error {

	if _, err := store.collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return err
	}

	for _, attribute := range []Attribute{
		FirstNameAttribute,
		LastNameAttribute,
		BirthDateAttribute,
	} {
		err := store.Create(ctx, &attribute)
		if err == nil {
			return nil
		}
		if mongo.IsDuplicateKeyError(err) {
			if err := store.Update(ctx, &attribute); err != nil {
				return err
			}
		}
	}
	return nil
}
