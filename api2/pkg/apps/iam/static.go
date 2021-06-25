package iam

import (
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
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
	PartyTypeIDs: []string{
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
	PartyTypeIDs: []string{
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
	PartyTypeIDs: []string{
		partytypes.IndividualPartyType.ID,
	},
}

var EMailAttribute = Attribute{
	ID:   "0ca7fa2b-982b-4fa5-85be-a6ebee8d4912",
	Name: "email",
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Email",
			ShortFormulation: "Email",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		partytypes.IndividualPartyType.ID,
	},
}

// StaffRelationshipType represents the built-in Staff relationship type
var StaffRelationshipType = RelationshipType{
	ID:              "53478121-23af-4ed8-a367-2e0de6d60271",
	Name:            "staff",
	FirstPartyRole:  "Is working for",
	SecondPartyRole: "Has staff",
	Rules: []RelationshipTypeRule{
		{
			PartyTypeRule: &PartyTypeRule{
				FirstPartyTypeID:  IndividualPartyType.ID,
				SecondPartyTypeID: OrganizationPartyType.ID,
			},
		},
	},
}

var IndividualPartyType = PartyType{
	ID:        "a842e7cb-3777-423a-9478-f1348be3b4a5",
	Name:      "Individual",
	IsBuiltIn: true,
}

var HouseholdPartyType = PartyType{
	ID:        "d38a7085-7dff-4730-8be1-7c9d92a20cc3",
	Name:      "Household",
	IsBuiltIn: true,
}

var OrganizationPartyType = partytypes.PartyType{
	ID:        "09a7eef9-3f23-4c40-86f4-9b9440c56c6f",
	Name:      "Organization",
	IsBuiltIn: true,
}
