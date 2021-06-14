package partytypes

import (
	"context"
)

var BeneficiaryPartyType = PartyType{
	ID:        "a842e7cb-3777-423a-9478-f1348be3b4a5",
	Name:      "Beneficiary",
	IsBuiltIn: true,
}

var HouseholdPartyType = PartyType{
	ID:        "d38a7085-7dff-4730-8be1-7c9d92a20cc3",
	Name:      "Household",
	IsBuiltIn: true,
}

func Init(ctx context.Context, store *Store) error {
	for _, partyType := range []PartyType{
		BeneficiaryPartyType,
		HouseholdPartyType,
	} {
		if err := store.Create(ctx, &partyType); err != nil {
			return err
		}
	}
	return nil
}
