package relationshiptypes

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var HeadOfHouseholdRelationshipType = api.RelationshipType{
	ID:              "de887604-9ce9-4fdc-af6b-602091a17913",
	Name:            "headOfHousehold",
	FirstPartyRole:  "Is head of household of",
	SecondPartyRole: "Has for head of household",
	Rules: []api.RelationshipTypeRule{
		{
			PartyTypeRule: api.PartyTypeRule{
				FirstPartyType:  partytypes.IndividualPartyType.ID,
				SecondPartyType: partytypes.HouseholdPartyType.ID,
			},
		},
	},
}

var SpousalRelationshipType = api.RelationshipType{
	ID:              "76376c69-ce06-4e06-b603-44c145ddf399",
	Name:            "spousal",
	FirstPartyRole:  "Is spouse of",
	SecondPartyRole: "Is spouse of",
	Rules: []api.RelationshipTypeRule{
		{
			PartyTypeRule: api.PartyTypeRule{
				FirstPartyType:  partytypes.IndividualPartyType.ID,
				SecondPartyType: partytypes.IndividualPartyType.ID,
			},
		},
	},
}

var FilialRelationshipType = api.RelationshipType{
	ID:              "dcebef97-f666-4593-b97e-075ad1890385",
	Name:            "filial",
	FirstPartyRole:  "Is sibling of",
	SecondPartyRole: "Is sibling of",
	Rules: []api.RelationshipTypeRule{
		{
			PartyTypeRule: api.PartyTypeRule{
				FirstPartyType:  partytypes.IndividualPartyType.ID,
				SecondPartyType: partytypes.IndividualPartyType.ID,
			},
		},
	},
}

var ParentalRelationshipType = api.RelationshipType{
	ID:              "628b9d26-f85d-44cd-8bed-6c5f692b4494",
	Name:            "parental",
	FirstPartyRole:  "Is parent of",
	SecondPartyRole: "Is child of",
	Rules: []api.RelationshipTypeRule{
		{
			PartyTypeRule: api.PartyTypeRule{
				FirstPartyType:  partytypes.IndividualPartyType.ID,
				SecondPartyType: partytypes.IndividualPartyType.ID,
			},
		},
	},
}

func Init(ctx context.Context, store *Store) error {

	if _, err := store.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"id": 1,
		},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return err
	}

	for _, relationshipType := range []api.RelationshipType{
		HeadOfHouseholdRelationshipType,
		SpousalRelationshipType,
		FilialRelationshipType,
		ParentalRelationshipType,
	} {
		if err := store.Create(ctx, &relationshipType); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := store.Update(ctx, &relationshipType); err != nil {
				return err
			}
		}
	}
	return nil

}
