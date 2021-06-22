package individuals

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var FirstNameAttribute = attributes.Attribute{
	ID:   "8514da51-aad5-4fb4-a797-8bcc0c969b27",
	Name: "firstName",
	Translations: []attributes.AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "First Name",
			ShortFormulation: "First Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		partytypes.IndividualPartyType.ID,
	},
}

var LastNameAttribute = attributes.Attribute{
	ID:   "21079bbc-e04b-4fe8-897f-644d73af0d9e",
	Name: "lastName",
	Translations: []attributes.AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Last Name",
			ShortFormulation: "Last Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		partytypes.IndividualPartyType.ID,
	},
}

var BirthDateAttribute = attributes.Attribute{
	ID:   "87fe07d7-e6a7-4428-8086-3842b69f3665",
	Name: "birthDate",
	Translations: []attributes.AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Birth Date",
			ShortFormulation: "Birth Date",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		partytypes.IndividualPartyType.ID,
	},
}

var BuiltinIndividualAttributes = []attributes.Attribute{
	FirstNameAttribute,
	LastNameAttribute,
	BirthDateAttribute,
}

var mockIndividuals = []*Individual{
	{
		Party: &parties.Party{
			ID:           "0bde06f0-5416-4514-9c5a-794a2cc2f1b7",
			PartyTypeIDs: []string{},
			Attributes: map[string][]string{
				FirstNameAttribute.ID: {"John"},
				LastNameAttribute.ID:  {"Doe"},
			},
		},
	}, {
		Party: &parties.Party{
			ID:           "ab7a1620-f34e-4811-8534-853167ed7944",
			PartyTypeIDs: []string{},
			Attributes: map[string][]string{
				FirstNameAttribute.ID: {"Mary"},
				LastNameAttribute.ID:  {"Poppins"},
			},
		},
	}, {
		Party: &parties.Party{
			ID:           "40b30fb0-c392-4798-9400-bda3e5837867",
			PartyTypeIDs: []string{},
			Attributes: map[string][]string{
				FirstNameAttribute.ID: {"Bo"},
				LastNameAttribute.ID:  {"Diddley"},
			},
		},
	},
}

func SeedDatabase(ctx context.Context, store *Store) error {
	for _, individual := range mockIndividuals {
		if err := store.Create(ctx, individual); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := store.Upsert(ctx, individual); err != nil {
				return err
			}
		}
	}
	return nil
}

func Init(ctx context.Context, store *attributes.Store, partiesStore *parties.Store) error {
	for _, attribute := range BuiltinIndividualAttributes {
		if err := store.Create(ctx, &attribute); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := store.Update(ctx, &attribute); err != nil {
				return err
			}
		}
	}

	// first name and last name full text index
	if _, err := partiesStore.Collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{
				"attributes." + FirstNameAttribute.ID, "text",
			},
			{
				"attributes." + LastNameAttribute.ID, "text",
			},
		},
	}); err != nil {
		return err
	}
	return nil
}
