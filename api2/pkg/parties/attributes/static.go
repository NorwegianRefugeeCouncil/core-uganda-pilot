package attributes

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
)

var FirstNameAttribute = api.Attribute{
	ID:   "8514da51-aad5-4fb4-a797-8bcc0c969b27",
	Name: "firstName",
	Translations: []api.AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "First Name",
			ShortFormulation: "First Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypes: []string{
		partytypes.BeneficiaryPartyType.ID,
	},
}

var LastNameAttribute = api.Attribute{
	ID:   "21079bbc-e04b-4fe8-897f-644d73af0d9e",
	Name: "lastName",
	Translations: []api.AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Last Name",
			ShortFormulation: "Last Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypes: []string{
		partytypes.BeneficiaryPartyType.ID,
	},
}

var BirthDateAttribute = api.Attribute{
	ID:   "87fe07d7-e6a7-4428-8086-3842b69f3665",
	Name: "birthDate",
	Translations: []api.AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Birth Date",
			ShortFormulation: "Birth Date",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypes: []string{
		partytypes.BeneficiaryPartyType.ID,
	},
}

func Init(ctx context.Context, store *Store) error {
	for _, attribute := range []api.Attribute{
		FirstNameAttribute,
		LastNameAttribute,
		BirthDateAttribute,
	} {
		if err := store.Create(ctx, &attribute); err != nil {
			return err
		}
	}
	return nil
}
