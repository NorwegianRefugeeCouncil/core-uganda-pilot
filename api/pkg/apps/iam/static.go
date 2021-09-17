package iam

import "github.com/nrc-no/core/pkg/form"

//Countries
var GlobalCountry = Country{
	ID:   "",
	Name: "Global",
}

var UgandaCountry = Country{
	ID:   "fc82a799-b4fc-4eda-81fc-f2710a0d27d8",
	Name: "Uganda",
}

var ColombiaCountry = Country{
	ID:   "d351395b-468c-4ceb-94d3-fa5f6338a5d3",
	Name: "Colombia",
}

// Global Individual Attributes -----------------------------------------------

var FullNameAttribute = Attribute{
	ID:         "8514da51-aad5-4fb4-a797-8bcc0c969b27",
	Name:       "fullName",
	CountryID:  GlobalCountry.ID,
	Type:       form.Text,
	Attributes: form.FormElementAttributes{Name: "fullName"},
	Validation: form.FormElementValidation{Required: true},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Full Name",
			ShortFormulation: "Full Name",
		},
	},
	IsPersonallyIdentifiableInfo: true,
	PartyTypeIDs: []string{
		IndividualPartyType.ID,
	},
}

var DisplayNameAttribute = Attribute{
	ID:         "21079bbc-e04b-4fe8-897f-644d73af0d9e",
	Name:       "displayName",
	CountryID:  GlobalCountry.ID,
	Type:       form.Text,
	Attributes: form.FormElementAttributes{Name: "displayName"},
	Validation: form.FormElementValidation{Required: true},
	Translations: []AttributeTranslation{
		{
			Locale:           "en",
			LongFormulation:  "Display Name",
			ShortFormulation: "Display Name",
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
	CountryID:  GlobalCountry.ID,
	Type:       form.Date,
	Attributes: form.FormElementAttributes{Name: "birthDate"},
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
	CountryID:  GlobalCountry.ID,
	Type:       form.Email,
	Attributes: form.FormElementAttributes{Name: "email"},
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

var DisplacementStatusAttribute = Attribute{
	ID:        "d1d824b2-d163-43ff-bc0a-527bd86b79bb",
	Name:      "displacementStatus",
	CountryID: GlobalCountry.ID,
	Type:      form.Dropdown,
	Attributes: form.FormElementAttributes{
		Name:    "displacementStatus",
		Options: []string{"Refugee", "Internally displaced person", "Host community", "Other"},
	},
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
	ID:        "b43f630c-2eb6-4629-af89-44ded61f7f3e",
	Name:      "gender",
	CountryID: GlobalCountry.ID,
	Type:      form.Dropdown,
	Attributes: form.FormElementAttributes{
		Name:    "gender",
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

var ConsentToNrcDataUseAttribute = Attribute{
	ID:        "8463d701-f964-4454-b8b2-efc202e8007d",
	Name:      "consent_to_nrc_data_use",
	CountryID: GlobalCountry.ID,
	Type:      form.Checkbox,
	Attributes: form.FormElementAttributes{
		Name: "consent_to_nrc_data_use",
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
	ID:         "1ac8cf17-49f3-4281-b9c9-6fd6036229c2",
	Name:       "consentToNrcDataUseProof",
	CountryID:  GlobalCountry.ID,
	Type:       form.URL,
	Attributes: form.FormElementAttributes{Name: "consentToNrcDataUseProof"},
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
	ID:        "0ab6fd31-fa0e-4d53-b236-94bce6f67d4b",
	Name:      "anonymous",
	CountryID: GlobalCountry.ID,
	Type:      form.Checkbox,
	Attributes: form.FormElementAttributes{
		Name: "anonymous",
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
	ID:        "24be4f47-ba00-405a-9bc5-c6fe58ecd80c",
	Name:      "minor",
	CountryID: GlobalCountry.ID,
	Type:      form.Checkbox,
	Attributes: form.FormElementAttributes{
		Name: "minor",
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
	ID:        "ae56b1fd-21f6-480a-9184-091a7093d8b8",
	Name:      "protectionConcerns",
	CountryID: GlobalCountry.ID,
	Type:      form.Checkbox,
	Attributes: form.FormElementAttributes{
		Name: "protectionConcerns",
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
	ID:        "cb51b2e8-27da-4375-b85f-c5c107f5d2b4",
	Name:      "physicalImpairment",
	CountryID: GlobalCountry.ID,
	Type:      form.Checkbox,
	Attributes: form.FormElementAttributes{
		Name: "physicalImpairment",
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
	ID:        "98def70b-ee72-40eb-aed1-5a834bf8f579",
	Name:      "physicalImpairmentIntensity",
	CountryID: GlobalCountry.ID,
	Type:      form.Dropdown,
	Attributes: form.FormElementAttributes{
		Name:    "physicalImpairmentIntensity",
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
	ID:        "972c0d7f-8fa9-436d-95ab-6773070bc451",
	Name:      "sensoryImpairment",
	CountryID: GlobalCountry.ID,
	Type:      form.Checkbox,
	Attributes: form.FormElementAttributes{
		Name: "sensoryImpairment",
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
	ID:        "b1e6cfac-a8b9-4a0d-a5c7-f164fde99bcc",
	Name:      "sensoryImpairmentIntensity",
	CountryID: GlobalCountry.ID,
	Type:      form.Dropdown,
	Attributes: form.FormElementAttributes{
		Name:    "sensoryImpairmentIntensity",
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

var MentalImpairmentAttribute = Attribute{
	ID:        "41b7eb87-6488-47e3-a4b0-1422c039d0c7",
	Name:      "mentalImpairment",
	CountryID: GlobalCountry.ID,
	Type:      form.Checkbox,
	Attributes: form.FormElementAttributes{
		Name: "mentalImpairment",
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
	ID:        "9983188b-4f43-4cd5-a972-fde3a08f4810",
	Name:      "mentalImpairmentIntensity",
	CountryID: GlobalCountry.ID,
	Type:      form.Dropdown,
	Attributes: form.FormElementAttributes{
		Name:    "mentalImpairmentIntensity",
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

// Global Individual Contact Attributes ---------------------------------------

var (
	PrimaryPhoneNumberAttribute = Attribute{
		ID:        "8eae83a8-cbc7-4ab2-a21f-d57cb3bb29ff",
		Name:      "primaryPhoneNumber",
		CountryID: GlobalCountry.ID,
		Type:      form.Phone,
		Attributes: form.FormElementAttributes{
			Name: "primaryPhoneNumber",
		},
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

	SecondaryPhoneNumberAttribute = Attribute{
		ID:        "1f3016af-ab39-422a-beb8-904b68a1619e",
		Name:      "secondaryPhoneNumber",
		CountryID: GlobalCountry.ID,
		Type:      form.Phone,
		Attributes: form.FormElementAttributes{
			Name: "secondaryPhoneNumber",
		},
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

	TertiaryPhoneNumberAttribute = Attribute{
		ID:        "4a0ba072-66a5-403f-bea1-35e9427659fb",
		Name:      "tertiaryPhoneNumber",
		CountryID: GlobalCountry.ID,
		Type:      form.Phone,
		Attributes: form.FormElementAttributes{
			Name: "tertiaryPhoneNumber",
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Tertiary phone number",
				ShortFormulation: "Tertiary phone number",
			},
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
)

// Uganda Individual Attributes -----------------------------------------------

var (
	UGIdentificationDateAttribute = Attribute{
		ID:        "c84b8b93-b974-4bec-b9f7-d437446b24a7",
		Name:      "ugIdentificationDate",
		CountryID: UgandaCountry.ID,
		Type:      form.Date,
		Attributes: form.FormElementAttributes{
			Label:       "Date of Identification",
			Name:        "ugIdentificationDate",
			Description: "Date of first interaction with NRC",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Date of Identification",
				ShortFormulation: "Date of Identification",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGIdentificationLocationAttribute = Attribute{
		ID:        "06680252-1a1f-4c9d-85dd-56feef20019d",
		Name:      "ugIdentificationLocation",
		CountryID: UgandaCountry.ID,
		Type:      form.Dropdown,
		Attributes: form.FormElementAttributes{
			Label:       "Location of Identification",
			Name:        "ugIdentificationLocation",
			Description: "",
			Options:     []string{"Kabusu Access Center", "Nsambya Access Center", "Kisenyi ICLA Center", "Lukuli ICLA Center", "Kawempe ICLA Center", "Ndejje ICLA Center", "Mengo Field Office", "Community (Specify location)", "Home Visit", "Phone", "Other (Specify)"},
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Location of Identification",
				ShortFormulation: "Location of Identification",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGIdentificationSourceAttribute = Attribute{
		ID:        "a131a0fb-0270-4feb-8fc9-46e7dd6b5acb",
		Name:      "ugIdentificationSource",
		CountryID: UgandaCountry.ID,
		Type:      form.Dropdown,
		Attributes: form.FormElementAttributes{
			Label:       "Source of Identification",
			Name:        "ugIdentificationSource",
			Description: "",
			Options:     []string{"Walk-in Center", "FFRM Referral", "Internal Referral (Other – Specify)", "ICLA Outreach Team", "External Referral (Community Leader/Contact)", "External Referral (INGO/LNGO)", "External Referral (Other – Specify)", "Self (Telephone)", "Self (Email)", "Internal Referral (Other NRC Sector – Specify)", "CBP Outreach Team", "Other NRC Outreach Team (Specify)", "External Referral (UN Agency)", "External Referral (Government)", "Other – Specify"},
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Source of Identification",
				ShortFormulation: "Source of Identification",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGAdmin2Attribute = Attribute{
		ID:        "44dffbc4-7536-42b9-af84-32ea4e9ed493",
		Name:      "ugAdmin2",
		CountryID: UgandaCountry.ID,
		Type:      form.Dropdown,
		Attributes: form.FormElementAttributes{
			Label:       "District / Admin 2",
			Name:        "ugAdmin2",
			Description: "",
			Options:     []string{"ABIM", "ADJUMANI", "ALEBTONG", "AMOLATAR", "AMUDAT", "AMURIA", "AMURU", "APAC", "BUDAKA", "BUGIRI", "BUIKWE", "BUKOMANSIMBI", "BUKWO", "BULAMBULI", "BULIISA", "BUNDIBUGYO", "BUSHENYI", "BUYENDE", "DOKOLO", "BUTAMBALA", "HOIMA", "IGANGA", "KAABONG", "KABALE", "KABAROLE", "KALANGALA", "KALIRO", "KALUNGU", "KAMULI", "KANUNGU", "KAPCHORWA", "KATAKWI", "KAYUNGA", "SHEEMA", "KITGUM", "KOBOKO", "KOLE", "KOTIDO", "KISORO", "KWEEN", "LAMWO", "LIRA", "LUUKA", "LYANTONDE", "MANAFWA", "MASAKA", "MASINDI", "MAYUGE", "MBALE", "MBARARA", "MOROTO", "MOYO", "NAKAPIRIPIRIT", "NAKASEKE", "NAKASONGOLA", "NAMUTUMBA", "NAPAK", "NEBBI", "NGORA", "BUHWEJU", "NTOROKO", "MARACHA", "OTUKE", "OYAM", "PADER", "RUBIRIZI", "SIRONKO", "SOROTI", "WAKISO", "YUMBE", "ZOMBO", "ISINGIRO", "MITOOMA", "KYEGEGWA", "NTUNGAMO", "RUKUNGIRI", "KAMWENGE", "IBANDA", "KASESE", "KIRUHURA", "KYENJOJO", "MUBENDE", "GOMBA", "KIBOGA", "MPIGI", "KYANKWANZI", "KAKUMIRO", "NWOYA", "KIRYANDONGO", "SERERE", "OMORO", "ARUA", "LWENGO", "SEMBABULE", "RAKAI", "MITYANA", "LUWERO", "MUKONO", "KAMPALA", "BUVUMA", "JINJA", "NAMAYINGO", "BUSIA", "BUDUDA", "TORORO", "BUTALEJA", "BUKEDEA", "KUMI", "PALLISA", "KIBUKU", "KABERAMAIDO", "AGAGO", "KAGADI", "KIBAALE", "GULU", "RUBANDA"},
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "District / Admin 2",
				ShortFormulation: "District / Admin 2",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGAdmin3Attribute = Attribute{
		ID:        "a17ffa5e-5d62-44cd-b89f-438eeba128ac",
		Name:      "ugAdmin3",
		CountryID: UgandaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Subcounty / Admin 3",
			Name:        "ugAdmin3",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Subcounty / Admin 3",
				ShortFormulation: "Subcounty / Admin 3",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGAdmin4Attribute = Attribute{
		ID:        "f867c62a-dcd0-4778-9f4e-7309d044e905",
		Name:      "ugAdmin4",
		CountryID: UgandaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Parish / Admin 4",
			Name:        "ugAdmin4",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Parish / Admin 4",
				ShortFormulation: "Parish / Admin 4",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGAdmin5Attribute = Attribute{
		ID:        "f0b34ffc-3e15-4195-8e90-a3e1e4b3940c",
		Name:      "ugAdmin5",
		CountryID: UgandaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Village / Admin 5",
			Name:        "ugAdmin5",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Village / Admin 5",
				ShortFormulation: "Village / Admin 5",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGNationalityAttribute = Attribute{
		ID:        "76aab836-73a6-4a1e-9c17-04b8a4c25d8d",
		Name:      "ugNationality",
		Type:      form.Dropdown,
		CountryID: UgandaCountry.ID,
		Attributes: form.FormElementAttributes{
			Name:     "ugNationality",
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

	UGSpokenLanguagesAttribute = Attribute{
		ID:        "d041cba5-9486-4390-bc2b-ec7fb03d67ff",
		Name:      "ugSpokenLanguages",
		Type:      form.Text,
		CountryID: UgandaCountry.ID,
		Attributes: form.FormElementAttributes{
			Name: "ugSpokenLanguages",
		},
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

	UGPreferredLanguageAttribute = Attribute{
		ID:        "da27a6e8-abe3-48d5-bfd9-46033e476a09",
		Name:      "ugPreferredLanguage",
		Type:      form.Text,
		CountryID: UgandaCountry.ID,
		Attributes: form.FormElementAttributes{
			Name: "ugPreferredLanguage",
		},
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

	UGPhysicalAddressAttribute = Attribute{
		ID:        "ac2795e8-15a5-42a0-b11f-b9269ff2a309",
		Name:      "ugPhysicalAddress",
		CountryID: UgandaCountry.ID,
		Type:      form.Textarea,
		Attributes: form.FormElementAttributes{
			Name: "ugPhysicalAddress",
		},
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

	UGInstructionOnMakingContactAttribute = Attribute{
		ID:        "4d399cb3-6653-4a61-92eb-331f07e6c395",
		Name:      "ugInstructionOnMakingContact",
		CountryID: GlobalCountry.ID,
		Type:      form.Textarea,
		Attributes: form.FormElementAttributes{
			Name: "ugInstructionOnMakingContact",
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Instructions on contacting the beneficiary",
				ShortFormulation: "Instructions on contacting the beneficiary",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGCanInitiateContactAttribute = Attribute{
		ID:        "7476fef0-d116-4b94-b981-ac647e16203d",
		Name:      "ugCanInitiateContact",
		CountryID: GlobalCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name: "ugCanInitiateContact",
			CheckboxOptions: []form.CheckboxOption{
				{Label: "NRC can initiate contact with Beneficiary."},
			},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "NRC can initiate contact with Beneficiary.",
				ShortFormulation: "NRC can initiate contact with Beneficiary.",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGPreferredMeansOfContactAttribute = Attribute{
		ID:        "1e7f2db9-eb63-46ae-b6d5-5c171a9e2534",
		Name:      "ugPreferredMeansOfContact",
		CountryID: UgandaCountry.ID,
		Type:      form.Dropdown,
		Attributes: form.FormElementAttributes{
			Name:    "ugPreferredMeansOfContact",
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

	UGRequireAnInterpreterAttribute = Attribute{
		ID:        "9b6ae87d-8935-49aa-9e32-26e7445d1afc",
		Name:      "ugRequireAnInterpreter",
		CountryID: UgandaCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name:            "ugRequireAnInterpreter",
			CheckboxOptions: []form.CheckboxOption{{Label: "This beneficiary requires an interpreter."}},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "This beneficiary requires an interpreter.",
				ShortFormulation: "Requires an interpreter",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
)

// Colombia Individual Attributes -----------------------------------------------

var (
	COPrimaryNationalityAttribute = Attribute{
		ID:        "d1ee17c5-a7c5-486f-a1e9-be4ec6d65700",
		Name:      "coPrimaryNationality",
		Type:      form.Dropdown,
		CountryID: ColombiaCountry.ID,
		Attributes: form.FormElementAttributes{
			Label:    "Primary nationality",
			Name:     "coPrimaryNationality",
			Options:  []string{"Colombia", "Venezuela", "Ecuador", "Panama", "Costa Rica", "Honduras"},
			Multiple: false,
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Primary Nationality",
				ShortFormulation: "Primary Nationality",
			},
			{
				Locale:           "es",
				LongFormulation:  "Nacionalidad (1)",
				ShortFormulation: "Nacionalidad (1)",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COSecondaryNationalityAttribute = Attribute{
		ID:        "74f39024-a318-4c6a-bb07-dfe55679f78f",
		Name:      "coSecondaryNationality",
		Type:      form.Dropdown,
		CountryID: ColombiaCountry.ID,
		Attributes: form.FormElementAttributes{
			Label:    "Secondary nationality",
			Name:     "coSecondaryNationality",
			Options:  []string{"Colombia", "Venezuela", "Ecuador", "Panama", "Costa Rica", "Honduras"},
			Multiple: false,
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Secondary Nationality",
				ShortFormulation: "Secondary Nationality",
			},
			{
				Locale:           "es",
				LongFormulation:  "Nacionalidad (2)",
				ShortFormulation: "Nacionalidad (2)",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COMaritalStatusAttribute = Attribute{
		ID:        "8bf6b645-20c1-403b-93bc-c05bbc22f570",
		Name:      "coMaritalStatus",
		Type:      form.Dropdown,
		CountryID: ColombiaCountry.ID,
		Attributes: form.FormElementAttributes{
			Name:     "coMaritalStatus",
			Label:    "Marital status",
			Options:  []string{"Married", "Single", "Divorced", "Separated", "Widdowed"},
			Multiple: true,
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Marital status",
				ShortFormulation: "Marital status",
			},
			{
				Locale:           "es",
				LongFormulation:  "Estado civil",
				ShortFormulation: "Estado civil",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COBeneficiaryTypeAttribute = Attribute{
		ID:        "796e4eb0-56a7-46bb-b81a-9727e674f1f8",
		Name:      "coBeneficiaryType",
		Type:      form.Dropdown,
		CountryID: ColombiaCountry.ID,
		Attributes: form.FormElementAttributes{
			Label:    "Beneficiary type",
			Name:     "coBeneficiaryType",
			Options:  []string{"Student", "Teacher", "Community leader", "Civil servant"},
			Multiple: true,
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Beneficiary type",
				ShortFormulation: "Beneficiary type",
			},
			{
				Locale:           "es",
				LongFormulation:  "Tipo de beneficiario\\a",
				ShortFormulation: "Tipo de beneficiario\\a",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COEthnicityAttribute = Attribute{
		ID:        "fe26bc55-30b7-4c30-97f1-99e90a3367a8",
		Name:      "coEthnicity",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Ethnicity",
			Name:        "coEthnicity",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Ethnicity",
				ShortFormulation: "Ethnicity",
			},
			{
				Locale:           "es",
				LongFormulation:  "Etnia",
				ShortFormulation: "Etnia",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	CORegistrationDateAttribute = Attribute{
		ID:        "7623b9f3-c29e-479f-872f-bd008a37aca4",
		Name:      "coRegistrationDate",
		CountryID: ColombiaCountry.ID,
		Type:      form.Date,
		Attributes: form.FormElementAttributes{
			Label:       "Registration date",
			Name:        "coRegistrationDate",
			Description: "Date of registration with NRC",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Registration date",
				ShortFormulation: "Registration date",
			},
			{
				Locale:           "es",
				LongFormulation:  "Fecha Registro",
				ShortFormulation: "Fecha Registro",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	CORegistrationLocationAttribute = Attribute{
		ID:        "f5ea04e0-7073-45b3-aa9a-a08afaf503da",
		Name:      "coRegistrationLocation",
		CountryID: ColombiaCountry.ID,
		Type:      form.Dropdown,
		Attributes: form.FormElementAttributes{
			Label:       "Location of Registration",
			Name:        "coRegistrationLocation",
			Description: "",
			Options:     []string{"Viento Libre", "Other (Specify)"},
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Location of Registration",
				ShortFormulation: "Location of Registration",
			},
			{
				Locale:           "es",
				LongFormulation:  "Lugar Registro",
				ShortFormulation: "Lugar Registro",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COSourceOfIdentificationAttribute = Attribute{
		ID:        "533dd5a3-8ab1-4eb0-9e20-8f4c7b02b2e9",
		Name:      "coSourceOfIdentification",
		CountryID: ColombiaCountry.ID,
		Type:      form.Dropdown,
		Attributes: form.FormElementAttributes{
			Label:       "Source of Identification",
			Name:        "coSourceOfIdentification",
			Description: "",
			Options:     []string{"Route", "Shelter", "Protective Space", "Home", "Community"},
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Source of Identification",
				ShortFormulation: "Source of Identification",
			},
			{
				Locale:           "es",
				LongFormulation:  "Tipo de lugar de atención",
				ShortFormulation: "Tipo de lugar de atención",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COTypeOfSettlementAttribute = Attribute{
		ID:        "ac56561b-64e4-4d96-bbe8-813a0ed7060c",
		Name:      "coTypeOfSettlement",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Type of settlement",
			Name:        "coTypeOfSettlement",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Type of settlement",
				ShortFormulation: "Type of settlement",
			},
			{
				Locale:           "es",
				LongFormulation:  "Tipo de asentamiento",
				ShortFormulation: "Tipo de asentamiento",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COEmergencyCareAttribute = Attribute{
		ID:        "c425da4b-5af1-4dff-abab-058b1cf9b122",
		Name:      "coEmergencyCare",
		CountryID: ColombiaCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name: "coEmergencyCare",
			CheckboxOptions: []form.CheckboxOption{
				{Label: "Beneficiary requires emergency care"},
			},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Does the beneficiary require emergency care?",
				ShortFormulation: "Requires emergency care",
			},
			{
				Locale:           "es",
				LongFormulation:  "Atención en emergencia",
				ShortFormulation: "Atención en emergencia",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	CODurableSolutionsAttribute = Attribute{
		ID:        "68241403-dd90-4e26-8d30-70db03b92c95",
		Name:      "coDurableSolutions",
		CountryID: ColombiaCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name: "coDurableSolutions",
			CheckboxOptions: []form.CheckboxOption{
				{Label: "Response is a durable solution?"},
			},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Is the response a durable solution?",
				ShortFormulation: "Durable solution",
			},
			{
				Locale:           "es",
				LongFormulation:  "Soluciones Duraderas",
				ShortFormulation: "Soluciones Duraderas",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COHardToReachAttribute = Attribute{
		ID:        "0c327266-47fb-4557-b2fc-a6e394432254",
		Name:      "coHardToReach",
		CountryID: ColombiaCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name: "coHardToReach",
			CheckboxOptions: []form.CheckboxOption{
				{Label: "Is the beneficiary in a hard to reach location?"},
			},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Is the beneficiary in a hard to reach location?",
				ShortFormulation: "Hard to reach",
			},
			{
				Locale:           "es",
				LongFormulation:  "Zona de difícil acceso",
				ShortFormulation: "Zona de difícil acceso",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAttendedCovid19Attribute = Attribute{
		ID:        "d59241dc-384b-430d-a8f4-f7851ff28615",
		Name:      "coAttendedCovid19",
		CountryID: ColombiaCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name: "coAttendedCovid19",
			CheckboxOptions: []form.CheckboxOption{
				{Label: "Did the beneficiary take part in Covid19 emergency training?"},
			},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Did the beneficiary take part in Covid19 emergency training?",
				ShortFormulation: "Attended Covid19",
			},
			{
				Locale:           "es",
				LongFormulation:  "Atendido Emergencia COVID-19",
				ShortFormulation: "Atendido Emergencia COVID-19",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COIntroSourceAttribute = Attribute{
		ID:        "dc7f97a3-b927-438e-9bdd-4374ae09b63a",
		Name:      "coIntroSource",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Name:  "coIntroSource",
			Label: "How was the beneficiary introduced to NRC?",
		},
		Validation: form.FormElementValidation{Required: false},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Intro source",
				ShortFormulation: "Intro source",
			},
			{
				Locale:           "es",
				LongFormulation:  "Como te has enterado de los servicios de NRC?",
				ShortFormulation: "Como te has enterado de los servicios de NRC?",
			},
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin1Attribute = Attribute{
		ID:        "88a5c89a-9f09-4513-a8cb-f81190f9cc0c",
		Name:      "coAdmin1",
		CountryID: ColombiaCountry.ID,
		Type:      form.Dropdown,
		Attributes: form.FormElementAttributes{
			Label:       "Country / Admin 1",
			Name:        "coAdmin1",
			Description: "",
			Options:     []string{"Colombia", "Venezuela", "Ecuador", "Panama", "Costa Rica", "Honduras"},
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Country / Admin 1",
				ShortFormulation: "Country / Admin 1",
			},
			{
				Locale:           "es",
				LongFormulation:  "País / Admin 1",
				ShortFormulation: "País / Admin 1",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin2Attribute = Attribute{
		ID:        "491d0ca0-0b63-4860-8e38-8139fcdccf51",
		Name:      "coAdmin2",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "District / Admin 2",
			Name:        "coAdmin2",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "District / Admin 2",
				ShortFormulation: "District / Admin 2",
			},
			{
				Locale:           "es",
				LongFormulation:  "Departamento / Admin 2",
				ShortFormulation: "Departamento / Admin 2",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin3Attribute = Attribute{
		ID:        "8e69cfdf-935e-43cb-81a0-79ebdda742ec",
		Name:      "coAdmin3",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Subcounty / Admin 3",
			Name:        "coAdmin3",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Subcounty / Admin 3",
				ShortFormulation: "Subcounty / Admin 3",
			},
			{
				Locale:           "es",
				LongFormulation:  "Municipio / Admin 3",
				ShortFormulation: "Municipio / Admin 3",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin4Attribute = Attribute{
		ID:        "cb132ade-f379-42a8-88b0-6c08b375e086",
		Name:      "coAdmin4",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Parish / Admin 4",
			Name:        "coAdmin4",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Parish / Admin 4",
				ShortFormulation: "Parish / Admin 4",
			},
			{
				Locale:           "es",
				LongFormulation:  "Comuna o Corregimiento / Admin 4",
				ShortFormulation: "Comuna o Corregimiento / Admin 4",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin5Attribute = Attribute{
		ID:        "faf65cc6-f5eb-4d18-91ca-00bbd3a3ab8e",
		Name:      "coAdmin5",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Village / Admin 5",
			Name:        "coAdmin5",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Village / Admin 5",
				ShortFormulation: "Village / Admin 5",
			},
			{
				Locale:           "es",
				LongFormulation:  "Barrio o Vereda / Admin 5",
				ShortFormulation: "Barrio o Vereda / Admin 5",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COJobOrEnterpriseAttribute = Attribute{
		ID:        "dda85258-9ce2-41e3-a7f8-21b837d65a25",
		Name:      "coJobOrEnterprise",
		CountryID: ColombiaCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name: "coJobOrEnterprise",
			CheckboxOptions: []form.CheckboxOption{
				{Label: "Do you have a job or enterprise?"},
			},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Do you have a job or enterprise?",
				ShortFormulation: "Do you have a job or enterprise?",
			},
			{
				Locale:           "es",
				LongFormulation:  "Ustedes tiene empleo o emprendimiento?",
				ShortFormulation: "Ustedes tiene empleo o emprendimiento?",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COTypeOfEnterpriseAttribute = Attribute{
		ID:        "94a9b0d8-8eb7-4165-ae3f-fcf7279da537",
		Name:      "coTypeOfEnterprise",
		Type:      form.Dropdown,
		CountryID: ColombiaCountry.ID,
		Attributes: form.FormElementAttributes{
			Label:    "Type of enterprise",
			Name:     "coTypeOfEnterprise",
			Options:  []string{"commerce", "production", "service", "agriculture"},
			Multiple: true,
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Type of enterprise",
				ShortFormulation: "Type of enterprise",
			},
			{
				Locale:           "es",
				LongFormulation:  "Tipo de emprendimiento",
				ShortFormulation: "Tipo de emprendimiento",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COTimeInBusinessAttribute = Attribute{
		ID:        "31e6f25b-d0a8-47c5-8161-0fdfdb39d430",
		Name:      "coTimeInBusiness",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Length of time in business",
			Name:        "coTimeInBusiness",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Length of time in business",
				ShortFormulation: "Length of time in business",
			},
			{
				Locale:           "es",
				LongFormulation:  "Tiempo del emprendimiento",
				ShortFormulation: "Tiempo del emprendimiento",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COTypeOfEmploymentAttribute = Attribute{
		ID:        "fd992e38-2ac1-41ed-9efb-1962174a6438",
		Name:      "coTypeOfEmployment",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Type of employment (type of contract)",
			Name:        "coTypeOfEmployment",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Type of employment (type of contract)",
				ShortFormulation: "Type of employment (type of contract)",
			},
			{
				Locale:           "es",
				LongFormulation:  "Tipo de empleo (modalidad de contrato)",
				ShortFormulation: "Tipo de empleo (modalidad de contrato)",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COFormsOfIncomeGenerationAttribute = Attribute{
		ID:        "0ac65773-d8ef-4f63-907b-9761e1630be8",
		Name:      "coFormsOfIncomeGeneration",
		CountryID: ColombiaCountry.ID,
		Type:      form.Textarea,
		Attributes: form.FormElementAttributes{
			Label:       "Forms of income generation in the family",
			Name:        "coFormsOfIncomeGeneration",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Forms of income generation in the family",
				ShortFormulation: "Forms of income generation in the family",
			},
			{
				Locale:           "es",
				LongFormulation:  "Formas de generacion de ingresos en la familia",
				ShortFormulation: "Formas de generacion de ingresos en la familia",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COLegalRepresentativeNameAttribue = Attribute{
		ID:        "48a22db2-e97d-44db-869b-0a192697c781",
		Name:      "coLegalRepresentativeName",
		CountryID: ColombiaCountry.ID,
		Type:      form.Text,
		Attributes: form.FormElementAttributes{
			Label:       "Name and surname of the legal representative",
			Name:        "coLegalRepresentativeName",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Name and surname of the legal representative",
				ShortFormulation: "Name and surname of the legal representative",
			},
			{
				Locale:           "es",
				LongFormulation:  "Nombre y apellido del representante juridico",
				ShortFormulation: "Nombre y apellido del representante juridico",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COLegalRepresentativeAdditionalInfoAttribute = Attribute{
		ID:        "f82b59f4-3610-4869-91a3-13308361d153",
		Name:      "coLegalRepresentativeAdditionalInfo",
		CountryID: ColombiaCountry.ID,
		Type:      form.Textarea,
		Attributes: form.FormElementAttributes{
			Label:       "Additional information about the legal representative",
			Name:        "coLegalRepresentativeAdditionalInfo",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Additional information about the legal representative",
				ShortFormulation: "Additional information about the legal representative",
			},
			{
				Locale:           "es",
				LongFormulation:  "Información adicional sobre el representante jurídico",
				ShortFormulation: "Información adicional sobre el representante jurídico",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COReasonsForRepresentationAttribute = Attribute{
		ID:        "42673d92-48fb-4426-b13f-104d2625a8ff",
		Name:      "coReasonsForRepresentation",
		CountryID: ColombiaCountry.ID,
		Type:      form.Textarea,
		Attributes: form.FormElementAttributes{
			Label:       "Reasons for representation",
			Name:        "coReasonsForRepresentation",
			Description: "",
		},
		Validation: form.FormElementValidation{Required: true},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Reasons for representation",
				ShortFormulation: "Reasons for representation",
			},
			{
				Locale:           "es",
				LongFormulation:  "Razones para representar",
				ShortFormulation: "Razones para representar",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COGuardianshipIsLegalAttribute = Attribute{
		ID:        "5af889e2-c3c9-4ceb-a6a4-d6fff7aa5747",
		Name:      "coGuardianshipIsLegal",
		CountryID: ColombiaCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name: "coGuardianshipIsLegal",
			CheckboxOptions: []form.CheckboxOption{
				{Label: "Is the guardianship legal according to national legislation?"},
			},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Is the guardianship legal according to national legislation?",
				ShortFormulation: "Is the guardianship legal according to national legislation?",
			},
			{
				Locale:           "es",
				LongFormulation:  "¿La tutela es legal según la legislación nacional?",
				ShortFormulation: "¿La tutela es legal según la legislación nacional?",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAbleToGiveLegalConsentAttribute = Attribute{
		ID:        "6d96c1e3-9a3c-40c8-93ae-36636ced0b1a",
		Name:      "coAbleToGiveLegalConsent",
		CountryID: ColombiaCountry.ID,
		Type:      form.Checkbox,
		Attributes: form.FormElementAttributes{
			Name: "coAbleToGiveLegalConsent",
			CheckboxOptions: []form.CheckboxOption{
				{Label: "Is the person able to give legal consent?"},
			},
		},
		Translations: []AttributeTranslation{
			{
				Locale:           "en",
				LongFormulation:  "Is the person able to give legal consent?",
				ShortFormulation: "Is the person able to give legal consent?",
			},
			{
				Locale:           "es",
				LongFormulation:  "La persona puede dar su consentimiento de forma legal?",
				ShortFormulation: "La persona puede dar su consentimiento de forma legal?",
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
)

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

var TeamNameAttribute = Attribute{
	ID:                           "18f410a3-6fde-45ce-80c7-fc5d92b85870",
	Name:                         "teamName",
	CountryID:                    GlobalCountry.ID,
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
	CountryID:                    GlobalCountry.ID,
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
