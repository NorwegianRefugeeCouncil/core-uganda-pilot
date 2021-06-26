package seed

import (
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
)

var mockIndividuals = []*iam.Individual{
	{
		Party: &iam.Party{
			ID:           "0bde06f0-5416-4514-9c5a-794a2cc2f1b7",
			PartyTypeIDs: []string{},
			Attributes: map[string][]string{
				iam.FirstNameAttribute.ID: {"John"},
				iam.LastNameAttribute.ID:  {"Doe"},
			},
		},
	}, {
		Party: &iam.Party{
			ID:           "ab7a1620-f34e-4811-8534-853167ed7944",
			PartyTypeIDs: []string{},
			Attributes: map[string][]string{
				iam.FirstNameAttribute.ID: {"Mary"},
				iam.LastNameAttribute.ID:  {"Poppins"},
			},
		},
	}, {
		Party: &iam.Party{
			ID:           "40b30fb0-c392-4798-9400-bda3e5837867",
			PartyTypeIDs: []string{},
			Attributes: map[string][]string{
				iam.FirstNameAttribute.ID: {"Bo"},
				iam.LastNameAttribute.ID:  {"Diddley"},
			},
		},
	},
}
