package iam

import "github.com/nrc-no/core/pkg/form"

var FirstNameAttribute = Attribute{
	ID:         "8514da51-aad5-4fb4-a797-8bcc0c969b27",
	Name:       "firstName",
	Type:       form.Text,
	Validation: form.FormElementValidation{Required: true},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "First Name",
			ShortFormulation: "First Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var LastNameAttribute = Attribute{
	ID:         "21079bbc-e04b-4fe8-897f-644d73af0d9e",
	Name:       "lastName",
	Type:       form.Text,
	Validation: form.FormElementValidation{Required: true},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Last Name",
			ShortFormulation: "Last Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var BirthDateAttribute = Attribute{
	ID:         "87fe07d7-e6a7-4428-8086-3842b69f3665",
	Name:       "birthDate",
	Type:       form.Date,
	Validation: form.FormElementValidation{Required: true},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Birth Date",
			ShortFormulation: "Birth Date",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var EMailAttribute = Attribute{
	ID:         "0ca7fa2b-982b-4fa5-85be-a6ebee8d4912",
	Name:       "email",
	Type:       form.Email,
	Validation: form.FormElementValidation{Required: true},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Email",
			ShortFormulation: "Email",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// ---------------------------------------------------------------------------

// Customisation for Uganda Demo

// TODO: This should be a dropdown attribute
// Values
// - Refugee
// - Internally-displaced person
// - Host community
var DisplacementStatusAttribute = Attribute{
	ID:         "d1d824b2-d163-43ff-bc0a-527bd86b79bb",
	Name:       "displacementStatus",
	Type:       form.Text,
	Validation: form.FormElementValidation{Required: true},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Displacement Status",
			ShortFormulation: "Displacement Status",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var GenderAttribute = Attribute{
	ID:   "b43f630c-2eb6-4629-af89-44ded61f7f3e",
	Name: "gender",
	Type: form.Dropdown,
	Attributes: form.FormElementAttributes{
		Options: []string{"Male", "Female"},
	},
	Validation: form.FormElementValidation{Required: true},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Gender",
			ShortFormulation: "Gender",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// TODO: This should be replaced with the OIDC consent mechanism
// This was previously mentioned by @ludydoo, and this attribute
// should be considered only for demo purposes!
// Also evaluate whether the proof attribute is still needed if
// using OIDC consent
var ConsentToNrcDataUseAttribute = Attribute{
	ID:   "8463d701-f964-4454-b8b2-efc202e8007d",
	Name: "consent_to_nrc_data_use",
	Type: form.Checkbox,
	Attributes: form.FormElementAttributes{
		CheckboxOptions: []form.CheckboxOption{
			{
				// FIXME Translation and checkbox option have no overlap
				Label: "Has the beneficiary consented to NRC using their data?",
				// Required: true,
			},
		},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Has the beneficiary consented to NRC using their data",
			ShortFormulation: "Consented to NRC using their data",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var ConsentToNrcDataUseProofAttribute = Attribute{
	ID:   "1ac8cf17-49f3-4281-b9c9-6fd6036229c2",
	Name: "consentToNrcDataUseProof",
	Type: form.File,
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Link to proof of beneficiary consent",
			ShortFormulation: "Consent proof",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var AnonymousAttribute = Attribute{
	ID:   "0ab6fd31-fa0e-4d53-b236-94bce6f67d4b",
	Name: "anonymous",
	Type: form.Checkbox,
	Attributes: form.FormElementAttributes{
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Beneficiary prefers to remain anonymous."},
		},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Beneficiary prefers to remain anonymous",
			ShortFormulation: "Anonymous",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var MinorAttribute = Attribute{
	ID:   "24be4f47-ba00-405a-9bc5-c6fe58ecd80c",
	Name: "minor",
	Type: form.Checkbox,
	Attributes: form.FormElementAttributes{
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Is the beneficiary a minor?"},
		},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Is the beneficiary a minor",
			ShortFormulation: "Minor",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var ProtectionConcernsAttribute = Attribute{
	ID:   "ae56b1fd-21f6-480a-9184-091a7093d8b8",
	Name: "protectionConcerns",
	Type: form.Checkbox,
	Attributes: form.FormElementAttributes{
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Beneficiary presents protection concerns"},
		},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Beneficiary presents protection concerns",
			ShortFormulation: "Protection concerns",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var PhysicalImpairmentAttribute = Attribute{
	ID:   "cb51b2e8-27da-4375-b85f-c5c107f5d2b4",
	Name: "physicalImpairment",
	Type: form.Checkbox,
	Attributes: form.FormElementAttributes{
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Would you say you experience some form of physical impairment?"},
		},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Would you say you experience some form of physical impairment?",
			ShortFormulation: "Experiences physical impairment",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var PhysicalImpairmentIntensityAttribute = Attribute{
	ID:   "98def70b-ee72-40eb-aed1-5a834bf8f579",
	Name: "physicalImpairmentIntensity",
	Type: form.Dropdown,
	Attributes: form.FormElementAttributes{
		Options: []string{"Moderate", "Severe"},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "How would you define the intensity of the physical impairment?",
			ShortFormulation: "Physical impairment intensity",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var SensoryImpairmentAttribute = Attribute{
	ID:   "972c0d7f-8fa9-436d-95ab-6773070bc451",
	Name: "sensoryImpairment",
	Type: form.Checkbox,
	Attributes: form.FormElementAttributes{
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Would you say you experience some form of sensory impairment?"},
		},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Would you say you experience some form of sensory impairment?",
			ShortFormulation: "Experiences sensory impairment",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var SensoryImpairmentIntensityAttribute = Attribute{
	ID:   "b1e6cfac-a8b9-4a0d-a5c7-f164fde99bcc",
	Name: "sensoryImpairmentIntensity",
	Type: form.Dropdown,
	Attributes: form.FormElementAttributes{
		Options: []string{"Moderate", "Severe"},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "How would you define the intensity of the sensory impairment?",
			ShortFormulation: "Sensory impairment intensity",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// TODO: Should be a checkbox/boolean attribute
var MentalImpairmentAttribute = Attribute{
	ID:   "41b7eb87-6488-47e3-a4b0-1422c039d0c7",
	Name: "mentalImpairment",
	Type: form.Checkbox,
	Attributes: form.FormElementAttributes{
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Would you say you experience some form of mental impairment?"},
		},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Would you say you experience some form of mental impairment?",
			ShortFormulation: "Experiences mental impairment",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var MentalImpairmentIntensityAttribute = Attribute{
	ID:   "9983188b-4f43-4cd5-a972-fde3a08f4810",
	Name: "sensoryImpairmentIntensity",
	Type: form.Dropdown,
	Attributes: form.FormElementAttributes{
		Options: []string{"Moderate", "Severe"},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "How would you define the intensity of the mental impairment?",
			ShortFormulation: "Mental impairment intensity",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var NationalityAttribute = Attribute{
	ID:   "76aab836-73a6-4a1e-9c17-04b8a4c25d8d",
	Name: "nationality",
	Type: form.Dropdown,
	Attributes: form.FormElementAttributes{
		Options:  []string{"Uganda", "Kenya", "Tanzania", "Rwanda", "Burundi", "Democratic Republic of Congo", "South Sudan", "Sudan", "Somalia", "Ethiopia"},
		Multiple: true,
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Nationality",
			ShortFormulation: "Nationality",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var SpokenLanguagesAttribute = Attribute{
	ID:   "d041cba5-9486-4390-bc2b-ec7fb03d67ff",
	Name: "spokenLanguages",
	Type: form.Text,
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "What languages does the beneficiary speak?",
			ShortFormulation: "Spoken languages",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var PreferredLanguageAttribute = Attribute{
	ID:   "da27a6e8-abe3-48d5-bfd9-46033e476a09",
	Name: "preferredLanguage",
	Type: form.Text,
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "What language does the beneficiary prefer for communication?",
			ShortFormulation: "Preferred language",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// TODO: Decide whether addresses should be their own entity
// This could allow beneficiaries to share addresses and reduce
// the work in maintaining data when a change needs to be made to
// an address in the app (1 update in 1 place)
var PhysicalAddressAttribute = Attribute{
	ID:   "ac2795e8-15a5-42a0-b11f-b9269ff2a309",
	Name: "physicalAddress",
	Type: form.Textarea,
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Physical address",
			ShortFormulation: "Physical address",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// TODO: Evaluate replacing primary + secondary numbers with an array type?
var PrimaryPhoneNumberAttribute = Attribute{
	ID:   "8eae83a8-cbc7-4ab2-a21f-d57cb3bb29ff",
	Name: "primaryPhoneNumber",
	Type: form.Phone,
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Primary phone number",
			ShortFormulation: "Primary phone number",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var SecondaryPhoneNumberAttribute = Attribute{
	ID:   "1f3016af-ab39-422a-beb8-904b68a1619e",
	Name: "secondaryPhoneNumber",
	Type: form.Phone,
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Secondary phone number",
			ShortFormulation: "Secondary phone number",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var PreferredMeansOfContactAttribute = Attribute{
	ID:   "1e7f2db9-eb63-46ae-b6d5-5c171a9e2534",
	Name: "preferredMeansOfContact",
	Type: form.Dropdown,
	Attributes: form.FormElementAttributes{
		Options: []string{"Phone Call", "Text message", "WhatsApp", "Signal", "Telegram", "Email", "Home visit"},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Preferred means of contact",
			ShortFormulation: "Preferred means of contact",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var RequireAnInterpreterAttribute = Attribute{
	ID:   "9b6ae87d-8935-49aa-9e32-26e7445d1afc",
	Name: "requireAnInterpreter",
	Type: form.Checkbox,
	Attributes: form.FormElementAttributes{
		CheckboxOptions: []form.CheckboxOption{{Label: "Does this beneficiary require an interpreter?"}},
	},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Does this beneficiary require an interpreter?",
			ShortFormulation: "Requires an interpreter",
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// ---------------------------------------------------------------------------

// StaffPartyType represents the built-in Staff relationship type
var StaffPartyType = PartyType{
	ID:        "53478121-23af-4ed8-a367-2e0de6d60271",
	Name:      "staff",
	IsBuiltIn: true,
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

var TeamPartyType = PartyType{
	ID:        "dacd6e08-3e3d-495b-8655-ea1d8e822cf3",
	Name:      "Team",
	IsBuiltIn: true,
}

var CountryPartyType = PartyType{
	ID:        "4954aaa1-98e4-480b-a542-3ffad12ca6cd",
	Name:      "Country",
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

var CountryNameAttribute = Attribute{
	ID:                           "e011d638-864b-496e-b3e5-af89d0278e1e",
	Name:                         "countryName",
	PartyTypeIDs:                 []string{CountryPartyType.ID},
	IsPersonallyIdentifiableInfo: false,
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			ShortFormulation: "Country name",
			LongFormulation:  "Country name",
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

var NationalityRelationshipType = RelationshipType{
	ID:              "4e9701db-7f5f-4536-a61f-b484997fe4c3",
	IsDirectional:   true,
	Name:            "teamNationality",
	FirstPartyRole:  "Is from country",
	SecondPartyRole: "Has team",
	Rules: []RelationshipTypeRule{
		{
			&PartyTypeRule{
				FirstPartyTypeID:  TeamPartyType.ID,
				SecondPartyTypeID: CountryPartyType.ID,
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
				FirstPartyTypeID:  IndividualPartyType.ID,
				SecondPartyTypeID: HouseholdPartyType.ID,
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
				FirstPartyTypeID:  IndividualPartyType.ID,
				SecondPartyTypeID: IndividualPartyType.ID,
			},
		},
	},
}

var SiblingRelationshipType = RelationshipType{
	ID:              "dcebef97-f666-4593-b97e-075ad1890385",
	IsDirectional:   false,
	Name:            "sibling",
	FirstPartyRole:  "Is sibling of",
	SecondPartyRole: "Is sibling of",
	Rules: []RelationshipTypeRule{
		{
			PartyTypeRule: &PartyTypeRule{
				FirstPartyTypeID:  IndividualPartyType.ID,
				SecondPartyTypeID: IndividualPartyType.ID,
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
				FirstPartyTypeID:  IndividualPartyType.ID,
				SecondPartyTypeID: IndividualPartyType.ID,
			},
		},
	},
}
