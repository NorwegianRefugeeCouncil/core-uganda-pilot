package iam

import "github.com/nrc-no/core/pkg/form"

//Countries
var GlobalCountry = Country{
	ID:   "36790d84-0bea-437c-b26e-bae1bcd2d1bc",
	Name: "Global",
}

var UgandaCountry = Country{
	ID:   "fc82a799-b4fc-4eda-81fc-f2710a0d27d8",
	Name: "Uganda",
}

//Individual Global Attributes
var FullNameAttribute = PartyAttributeDefinition{
	ID: "8514da51-aad5-4fb4-a797-8bcc0c969b27",
	FormControl: form.Control{
		Name:       "fullName",
		Type:       form.Text,
		Label:      "Full Name",
		Validation: form.ControlValidation{Required: true},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var DisplayNameAttribute = PartyAttributeDefinition{
	ID: "21079bbc-e04b-4fe8-897f-644d73af0d9e",
	FormControl: form.Control{
		Name:       "displayName",
		Type:       form.Text,
		Label:      "Display Name",
		Validation: form.ControlValidation{Required: true},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var BirthDateAttribute = PartyAttributeDefinition{
	ID: "87fe07d7-e6a7-4428-8086-3842b69f3665",
	FormControl: form.Control{
		Name:  "birthDate",
		Type:  form.Date,
		Label: "Birth Date",
		Validation: form.ControlValidation{
			Required: true,
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var EMailAttribute = PartyAttributeDefinition{
	ID: "0ca7fa2b-982b-4fa5-85be-a6ebee8d4912",
	FormControl: form.Control{
		Name:       "email",
		Type:       form.Email,
		Label:      "Email",
		Validation: form.ControlValidation{Required: true},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// Uganda Idividual Attributes
var IdentificationDateAttribute = PartyAttributeDefinition{
	ID: "c84b8b93-b974-4bec-b9f7-d437446b24a7",
	FormControl: form.Control{
		Name:        "identificationDate",
		Type:        form.Date,
		Label:       "Date of identification",
		Description: "Date of first interaction with NRC",
	},
	CountryID:                    UgandaCountry.ID,
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var IdentificationLocationAttribute = PartyAttributeDefinition{
	ID:        "06680252-1a1f-4c9d-85dd-56feef20019d",
	CountryID: UgandaCountry.ID,
	FormControl: form.Control{
		Name:  "identificationLocation",
		Type:  form.Dropdown,
		Label: "Location of Identification",
		Options: []string{
			"Kabusu Access Center",
			"Nsambya Access Center",
			"Kisenyi ICLA Center",
			"Lukuli ICLA Center",
			"Kawempe ICLA Center",
			"Ndejje ICLA Center",
			"Mengo Field Office",
			"Community (Specify location)",
			"Home Visit",
			"Phone",
			"Other (Specify)",
		},
		Validation: form.ControlValidation{
			Required: true,
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var IdentificationSourceAttribute = PartyAttributeDefinition{
	ID: "a131a0fb-0270-4feb-8fc9-46e7dd6b5acb",
	FormControl: form.Control{
		Name:       "identificationSource",
		Type:       form.Dropdown,
		Label:      "Source of Identification",
		Validation: form.ControlValidation{Required: true},
	},
	CountryID:                    UgandaCountry.ID,
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var Admin2Attribute = PartyAttributeDefinition{
	ID:        "44dffbc4-7536-42b9-af84-32ea4e9ed493",
	CountryID: UgandaCountry.ID,
	FormControl: form.Control{
		Name:    "admin2",
		Type:    form.Dropdown,
		Label:   "District / Admin 2",
		Options: []string{"ABIM", "ADJUMANI", "ALEBTONG", "AMOLATAR", "AMUDAT", "AMURIA", "AMURU", "APAC", "BUDAKA", "BUGIRI", "BUIKWE", "BUKOMANSIMBI", "BUKWO", "BULAMBULI", "BULIISA", "BUNDIBUGYO", "BUSHENYI", "BUYENDE", "DOKOLO", "BUTAMBALA", "HOIMA", "IGANGA", "KAABONG", "KABALE", "KABAROLE", "KALANGALA", "KALIRO", "KALUNGU", "KAMULI", "KANUNGU", "KAPCHORWA", "KATAKWI", "KAYUNGA", "SHEEMA", "KITGUM", "KOBOKO", "KOLE", "KOTIDO", "KISORO", "KWEEN", "LAMWO", "LIRA", "LUUKA", "LYANTONDE", "MANAFWA", "MASAKA", "MASINDI", "MAYUGE", "MBALE", "MBARARA", "MOROTO", "MOYO", "NAKAPIRIPIRIT", "NAKASEKE", "NAKASONGOLA", "NAMUTUMBA", "NAPAK", "NEBBI", "NGORA", "BUHWEJU", "NTOROKO", "MARACHA", "OTUKE", "OYAM", "PADER", "RUBIRIZI", "SIRONKO", "SOROTI", "WAKISO", "YUMBE", "ZOMBO", "ISINGIRO", "MITOOMA", "KYEGEGWA", "NTUNGAMO", "RUKUNGIRI", "KAMWENGE", "IBANDA", "KASESE", "KIRUHURA", "KYENJOJO", "MUBENDE", "GOMBA", "KIBOGA", "MPIGI", "KYANKWANZI", "KAKUMIRO", "NWOYA", "KIRYANDONGO", "SERERE", "OMORO", "ARUA", "LWENGO", "SEMBABULE", "RAKAI", "MITYANA", "LUWERO", "MUKONO", "KAMPALA", "BUVUMA", "JINJA", "NAMAYINGO", "BUSIA", "BUDUDA", "TORORO", "BUTALEJA", "BUKEDEA", "KUMI", "PALLISA", "KIBUKU", "KABERAMAIDO", "AGAGO", "KAGADI", "KIBAALE", "GULU", "RUBANDA"},
		Validation: form.ControlValidation{
			Required: true,
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var Admin3Attribute = PartyAttributeDefinition{
	ID: "a17ffa5e-5d62-44cd-b89f-438eeba128ac",
	FormControl: form.Control{
		Name:  "admin3",
		Type:  form.Text,
		Label: "Subcounty / Admin 3",
		Validation: form.ControlValidation{
			Required: true,
		},
	},
	CountryID:                    UgandaCountry.ID,
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var Admin4Attribute = PartyAttributeDefinition{
	ID:        "f867c62a-dcd0-4778-9f4e-7309d044e905",
	CountryID: UgandaCountry.ID,
	FormControl: form.Control{
		Name:        "admin4",
		Type:        form.Text,
		Label:       "Parish / Admin 4",
		Description: "",
		Validation: form.ControlValidation{
			Required: true,
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var Admin5Attribute = PartyAttributeDefinition{
	ID: "f0b34ffc-3e15-4195-8e90-a3e1e4b3940c",
	FormControl: form.Control{
		Name:  "admin5",
		Type:  form.Text,
		Label: "Village / Admin 5",
		Validation: form.ControlValidation{
			Required: true,
		},
	},
	CountryID:                    UgandaCountry.ID,
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// ---------------------------------------------------------------------------

var DisplacementStatusAttribute = PartyAttributeDefinition{
	ID: "d1d824b2-d163-43ff-bc0a-527bd86b79bb",
	FormControl: form.Control{
		Name:    "displacementStatus",
		Type:    form.Dropdown,
		Label:   "Displacement Status",
		Options: []string{"Refugee", "Internally displaced person", "Host community", "Other"},
		Validation: form.ControlValidation{
			Required: true,
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var GenderAttribute = PartyAttributeDefinition{
	ID: "b43f630c-2eb6-4629-af89-44ded61f7f3e",
	FormControl: form.Control{
		Name:       "gender",
		Type:       form.Dropdown,
		Label:      "Gender",
		Validation: form.ControlValidation{Required: true},
		Options:    []string{"Male", "Female"},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// TODO: This should be replaced with the OIDC consent mechanism
// COR-208
// This was previously mentioned by @ludydoo, and this attribute
// should be considered only for demo purposes!
// Also evaluate whether the proof attribute is still needed if
// using OIDC consent
var ConsentToNrcDataUseAttribute = PartyAttributeDefinition{
	ID: "8463d701-f964-4454-b8b2-efc202e8007d",
	FormControl: form.Control{
		Name: "consentToDataUse",
		Type: form.Checkbox,
		CheckboxOptions: []form.CheckboxOption{
			{
				Label: "Has the beneficiary consented to NRC using their data?",
			},
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var ConsentToNrcDataUseProofAttribute = PartyAttributeDefinition{
	ID: "1ac8cf17-49f3-4281-b9c9-6fd6036229c2",
	FormControl: form.Control{
		Name:  "consentLink",
		Type:  form.URL,
		Label: "Link to proof of beneficiary consent",
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var AnonymousAttribute = PartyAttributeDefinition{
	ID: "0ab6fd31-fa0e-4d53-b236-94bce6f67d4b",
	FormControl: form.Control{
		Name: "beneficiaryPrefersAnonymous",
		Type: form.Checkbox,
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Beneficiary prefers to remain anonymous."},
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var MinorAttribute = PartyAttributeDefinition{
	ID: "24be4f47-ba00-405a-9bc5-c6fe58ecd80c",
	FormControl: form.Control{
		Name: "beneficiaryIsMinor",
		Type: form.Checkbox,
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Is the beneficiary a minor?"},
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var ProtectionConcernsAttribute = PartyAttributeDefinition{
	ID: "ae56b1fd-21f6-480a-9184-091a7093d8b8",
	FormControl: form.Control{
		Name: "protectionConcerns",
		Type: form.Checkbox,
		CheckboxOptions: []form.CheckboxOption{
			{
				Label: "Beneficiary presents protection concerns",
			},
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var PhysicalImpairmentAttribute = PartyAttributeDefinition{
	ID: "cb51b2e8-27da-4375-b85f-c5c107f5d2b4",
	FormControl: form.Control{
		Name:  "physicalImpairment",
		Type:  form.Checkbox,
		Label: "Physical impairment",
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Would you say you experience some form of physical impairment?"},
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var PhysicalImpairmentIntensityAttribute = PartyAttributeDefinition{
	ID: "98def70b-ee72-40eb-aed1-5a834bf8f579",
	FormControl: form.Control{
		Name:        "physicalImpairmentIntensity",
		Type:        form.Dropdown,
		Label:       "Physical impairment intensity",
		Description: "How would you define the intensity of the physical impairment?",
		Options:     []string{"Moderate", "Severe"},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var SensoryImpairmentAttribute = PartyAttributeDefinition{
	ID: "972c0d7f-8fa9-436d-95ab-6773070bc451",
	FormControl: form.Control{
		Name:  "sensoryImpairment",
		Type:  form.Checkbox,
		Label: "Sensory impairement",
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Would you say you experience some form of sensory impairment?"},
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var SensoryImpairmentIntensityAttribute = PartyAttributeDefinition{
	ID: "b1e6cfac-a8b9-4a0d-a5c7-f164fde99bcc",
	FormControl: form.Control{
		Name:    "sensoryImpairmentIntensity",
		Type:    form.Dropdown,
		Label:   "Sensory Impairment Intensity",
		Options: []string{"Moderate", "Severe"},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var MentalImpairmentAttribute = PartyAttributeDefinition{
	ID: "41b7eb87-6488-47e3-a4b0-1422c039d0c7",
	FormControl: form.Control{
		Name:  "mentalImpairment",
		Type:  form.Checkbox,
		Label: "Would you say you experience some form of mental impairment?",
		CheckboxOptions: []form.CheckboxOption{
			{Label: "Would you say you experience some form of mental impairment?"},
		},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var MentalImpairmentIntensityAttribute = PartyAttributeDefinition{
	ID: "9983188b-4f43-4cd5-a972-fde3a08f4810",
	FormControl: form.Control{
		Name:        "mentalImpairmentIntensity",
		Type:        form.Dropdown,
		Label:       "Mental impairment intensity",
		Options:     []string{"Moderate", "Severe"},
		Description: "How would you define the intensity of the mental impairment?",
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var NationalityAttribute = PartyAttributeDefinition{
	ID: "76aab836-73a6-4a1e-9c17-04b8a4c25d8d",
	FormControl: form.Control{
		Name:     "nationality",
		Type:     form.Dropdown,
		Label:    "Nationality",
		Options:  []string{"Uganda", "Kenya", "Tanzania", "Rwanda", "Burundi", "Democratic Republic of Congo", "South Sudan", "Sudan", "Somalia", "Ethiopia"},
		Multiple: true,
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var SpokenLanguagesAttribute = PartyAttributeDefinition{
	ID: "d041cba5-9486-4390-bc2b-ec7fb03d67ff",
	FormControl: form.Control{
		Name:        "spokenLanguages",
		Label:       "Spoken languages",
		Description: "What languages does the beneficiary speak?",
		Type:        form.Text,
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var PreferredLanguageAttribute = PartyAttributeDefinition{
	ID: "da27a6e8-abe3-48d5-bfd9-46033e476a09",
	FormControl: form.Control{
		Name:  "preferredLanguage",
		Type:  form.Text,
		Label: "What language does the beneficiary prefer for communication?",
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
var PhysicalAddressAttribute = PartyAttributeDefinition{
	ID: "ac2795e8-15a5-42a0-b11f-b9269ff2a309",
	FormControl: form.Control{
		Name:  "physicalAddress",
		Type:  form.Textarea,
		Label: "Physical address",
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

// TODO: Evaluate replacing primary + secondary numbers with an array type?
var PrimaryPhoneNumberAttribute = PartyAttributeDefinition{
	ID: "8eae83a8-cbc7-4ab2-a21f-d57cb3bb29ff",
	FormControl: form.Control{
		Name:  "phonePrimary",
		Type:  form.Phone,
		Label: "Primary phone number",
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var SecondaryPhoneNumberAttribute = PartyAttributeDefinition{
	ID: "1f3016af-ab39-422a-beb8-904b68a1619e",
	FormControl: form.Control{
		Name:  "phoneSecondary",
		Type:  form.Phone,
		Label: "Secondary phone number",
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var PreferredMeansOfContactAttribute = PartyAttributeDefinition{
	ID: "1e7f2db9-eb63-46ae-b6d5-5c171a9e2534",
	FormControl: form.Control{
		Name:    "preferredMeansOfContact",
		Type:    form.Dropdown,
		Label:   "Preferred means of contact",
		Options: []string{"Phone Call", "Text message", "WhatsApp", "Signal", "Telegram", "Email", "Home visit"},
	},
	IsPersonallyIdentifiableInfo: false,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var RequireAnInterpreterAttribute = PartyAttributeDefinition{
	ID: "9b6ae87d-8935-49aa-9e32-26e7445d1afc",
	FormControl: form.Control{
		Name:            "requiresInterpreter",
		Type:            form.Checkbox,
		CheckboxOptions: []form.CheckboxOption{{Label: "This beneficiary requires an interpreter."}},
		Description:     "This beneficiary requires an interpreter.",
		Label:           "Requires an interpreter",
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

var BeneficiaryPartyType = PartyType{
	ID:        "09dbb93e-25c5-4cd5-a861-c2706efee0e0",
	Name:      "beneficiary",
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

var TeamNameAttribute = PartyAttributeDefinition{
	ID: "18f410a3-6fde-45ce-80c7-fc5d92b85870",
	FormControl: form.Control{
		Name:  "teamName",
		Type:  form.Text,
		Label: "Team name",
	},
	PartyTypeIDs:                 []string{TeamPartyType.ID},
	IsPersonallyIdentifiableInfo: false,
}

var CountryNameAttribute = PartyAttributeDefinition{
	ID: "e011d638-864b-496e-b3e5-af89d0278e1e",
	FormControl: form.Control{
		Name:  "countryName",
		Type:  form.Text,
		Label: "Country name",
	},
	PartyTypeIDs:                 []string{CountryPartyType.ID},
	IsPersonallyIdentifiableInfo: false,
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
