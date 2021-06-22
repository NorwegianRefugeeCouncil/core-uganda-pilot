package individuals

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"go.mongodb.org/mongo-driver/mongo"
)

var mockIndividuals = []*Individual{
	{
		Party: &parties.Party{
			ID:         "0bde06f0-5416-4514-9c5a-794a2cc2f1b7",
			PartyTypes: []string{},
			Attributes: map[string][]string{
				attributes.FirstNameAttribute.ID: {"John"},
				attributes.LastNameAttribute.ID:  {"Doe"},
			},
		},
	}, {
		Party: &parties.Party{
			ID:         "ab7a1620-f34e-4811-8534-853167ed7944",
			PartyTypes: []string{},
			Attributes: map[string][]string{
				attributes.FirstNameAttribute.ID: {"Mary"},
				attributes.LastNameAttribute.ID:  {"Poppins"},
			},
		},
	}, {
		Party: &parties.Party{
			ID:         "40b30fb0-c392-4798-9400-bda3e5837867",
			PartyTypes: []string{},
			Attributes: map[string][]string{
				attributes.FirstNameAttribute.ID: {"Bo"},
				attributes.LastNameAttribute.ID:  {"Diddley"},
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
