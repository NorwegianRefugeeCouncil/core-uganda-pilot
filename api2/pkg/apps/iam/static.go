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

var OrganizationPartyType = PartyType{
	ID:        "09a7eef9-3f23-4c40-86f4-9b9440c56c6f",
	Name:      "Organization",
	IsBuiltIn: true,
}

var TeamPartyType = PartyType{
	ID:        "dacd6e08-3e3d-495b-8655-ea1d8e822cf3",
	Name:      "Team",
	IsBuiltIn: true,
}

var TeamNameAttribute = Attribute{
	ID:                           "18f410a3-6fde-45ce-80c7-fc5d92b85870",
	Name:                         "teamName",
	PartyTypeIDs:                 []string{TeamPartyType.ID},
	IsPersonallyIdentifiableInfo: false,
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			ShortFormulation: "Team name",
			LongFormulation:  "Team name",
		},
	},
}

var MembershipRelationshipType = RelationshipType{
	ID:              "69fef57b-b37f-4803-a5fb-47e05282ac84",
	IsDirectional:   true,
	Name:            "teamMembership",
	FirstPartyRole:  "Is member of team",
	SecondPartyRole: "Has team member",
	Rules: []RelationshipTypeRule{
		{
			&PartyTypeRule{
				FirstPartyTypeID:  IndividualPartyType.ID,
				SecondPartyTypeID: TeamPartyType.ID,
			},
		},
	},
}

var HeadOfHouseholdRelationshipType = RelationshipType{
	ID:              "de887604-9ce9-4fdc-af6b-602091a17913",
	IsDirectional:   true,
	Name:            "headOfHousehold",
	FirstPartyRole:  "Is head of household of",
	SecondPartyRole: "Has for head of household",
	Rules: []RelationshipTypeRule{
		{
			PartyTypeRule: &PartyTypeRule{
				FirstPartyTypeID:  partytypes.IndividualPartyType.ID,
				SecondPartyTypeID: partytypes.HouseholdPartyType.ID,
			},
		},
	},
}

var SpousalRelationshipType = RelationshipType{
	ID:              "76376c69-ce06-4e06-b603-44c145ddf399",
	IsDirectional:   false,
	Name:            "spousal",
	FirstPartyRole:  "Is spouse of",
	SecondPartyRole: "Is spouse of",
	Rules: []RelationshipTypeRule{
		{
			PartyTypeRule: &PartyTypeRule{
				FirstPartyTypeID:  partytypes.IndividualPartyType.ID,
				SecondPartyTypeID: partytypes.IndividualPartyType.ID,
			},
		},
	},
}

var FilialRelationshipType = RelationshipType{
	ID:              "dcebef97-f666-4593-b97e-075ad1890385",
	IsDirectional:   false,
	Name:            "filial",
	FirstPartyRole:  "Is sibling of",
	SecondPartyRole: "Is sibling of",
	Rules: []RelationshipTypeRule{
		{
			PartyTypeRule: &PartyTypeRule{
				FirstPartyTypeID:  partytypes.IndividualPartyType.ID,
				SecondPartyTypeID: partytypes.IndividualPartyType.ID,
			},
		},
	},
}

var ParentalRelationshipType = RelationshipType{
	ID:              "628b9d26-f85d-44cd-8bed-6c5f692b4494",
	IsDirectional:   true,
	Name:            "parental",
	FirstPartyRole:  "Is parent of",
	SecondPartyRole: "Is child of",
	Rules: []RelationshipTypeRule{
		{
			PartyTypeRule: &PartyTypeRule{
				FirstPartyTypeID:  partytypes.IndividualPartyType.ID,
				SecondPartyTypeID: partytypes.IndividualPartyType.ID,
			},
		},
	},
}

var LegalNameAttribute = Attribute{
	ID:   "7afb0744-c764-4c5b-9dc6-b341d9b320b4",
	Name: "legalName",
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Legal Name",
			ShortFormulation: "Legal Name",
		},
	},
	PartyTypeIDs: []string{
		OrganizationPartyType.ID,
	},
}
