package beneficiaries

import (
	"context"
	api2 "github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/attributes"
	"github.com/nrc-no/core-kafka/pkg/parties/beneficiaries/api"
	"go.mongodb.org/mongo-driver/mongo"
)

var mockBeneficiaries = []*api.Beneficiary{
	{
		Party: &api2.Party{
			ID:         "0bde06f0-5416-4514-9c5a-794a2cc2f1b7",
			PartyTypes: []string{},
			Attributes: map[string][]string{
				attributes.FirstNameAttribute.ID: {"John"},
				attributes.LastNameAttribute.ID:  {"Doe"},
			},
		},
	}, {
		Party: &api2.Party{
			ID:         "ab7a1620-f34e-4811-8534-853167ed7944",
			PartyTypes: []string{},
			Attributes: map[string][]string{
				attributes.FirstNameAttribute.ID: {"Mary"},
				attributes.LastNameAttribute.ID:  {"Poppins"},
			},
		},
	}, {
		Party: &api2.Party{
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
	for _, beneficiary := range mockBeneficiaries {
		if err := store.Create(ctx, beneficiary); err != nil {
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
			if err := store.Upsert(ctx, beneficiary); err != nil {
				return err
			}
		}
	}
	return nil
}
