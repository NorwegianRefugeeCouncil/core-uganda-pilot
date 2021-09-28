package iam

import (
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/i18n"
)

//Countries

var (
	GlobalCountry = Country{
		ID:   "",
		Name: "Global",
	}
	UgandaCountry = Country{
		ID:   "fc82a799-b4fc-4eda-81fc-f2710a0d27d8",
		Name: "Uganda",
	}
	ColombiaCountry = Country{
		ID:   "d351395b-468c-4ceb-94d3-fa5f6338a5d3",
		Name: "Colombia",
	}
)

// Global Individual Attributes -----------------------------------------------

var (
	FullNameAttribute = PartyAttributeDefinition{
		ID:        "8514da51-aad5-4fb4-a797-8bcc0c969b27",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "fullName",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Full Name"},
				{"es", "Nombre completo"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	DisplayNameAttribute = PartyAttributeDefinition{
		ID:        "21079bbc-e04b-4fe8-897f-644d73af0d9e",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "displayName",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Display Name"},
				{"es", "Nombre para mostrar"},
			},

			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	BirthDateAttribute = PartyAttributeDefinition{
		ID:        "87fe07d7-e6a7-4428-8086-3842b69f3665",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "birthDate",
			Type: form.Date,
			Label: i18n.Strings{
				{"en", "Birth Date"},
				{"es", "Fecha de nacimiento"},
			},
			Validation: form.ControlValidation{
				Required: true,
			},
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	EMailAttribute = PartyAttributeDefinition{
		ID:        "0ca7fa2b-982b-4fa5-85be-a6ebee8d4912",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "email",
			Type: form.Email,
			Label: i18n.Strings{
				{"en", "Email"},
				{"es", "Email"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	DisplacementStatusAttribute = PartyAttributeDefinition{
		ID:        "d1d824b2-d163-43ff-bc0a-527bd86b79bb",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "displacementStatus",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Displacement Status"},
				{"es", "Estado de desplazamiento"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Refugee"},
					{"es", "Refugiada/o"},
				},
				{
					{"en", "Internally displaced person"},
					{"es", "Persona desplazada internamente"},
				},
				{
					{"en", "Host community"},
					{"es", "Comunidad anfitriona"},
				},
				{
					{"en", "Other"},
					{"es", "Otro"},
				},
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
	GenderAttribute = PartyAttributeDefinition{
		ID:        "b43f630c-2eb6-4629-af89-44ded61f7f3e",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "gender",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Gender"},
				{"es", "Género"},
			},
			Validation: form.ControlValidation{Required: true},
			Options: []i18n.Strings{
				{
					{"en", "Male"},
					{"es", "Masculino"},
				},
				{
					{"en", "Female"},
					{"es", "Femenino"},
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	ConsentToNrcDataUseAttribute = PartyAttributeDefinition{
		ID:        "8463d701-f964-4454-b8b2-efc202e8007d",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "consentToDataUse",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Has the beneficiary consented to NRC using their data?"},
						{"es", "¿El beneficiario ha dado su consentimiento para que NRC use sus datos?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	ConsentToNrcDataUseProofAttribute = PartyAttributeDefinition{
		ID:        "1ac8cf17-49f3-4281-b9c9-6fd6036229c2",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "consentToNrcDataUseProof",
			Type: form.URL,
			Label: i18n.Strings{
				{"en", "Link to proof of beneficiary consent"},
				{"es", "Enlace a prueba de consentimiento del beneficiario"},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	AnonymousAttribute = PartyAttributeDefinition{
		ID:        "0ab6fd31-fa0e-4d53-b236-94bce6f67d4b",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "anonymous",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Beneficiary prefers to remain anonymous."},
						{"es", "La/el beneficiaria/o prefiere permanecer anónima/o."},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	MinorAttribute = PartyAttributeDefinition{
		ID:        "24be4f47-ba00-405a-9bc5-c6fe58ecd80c",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "minor",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Is the beneficiary a minor?"},
						{"es", "¿La/el beneficiaria/o es menor de edad?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	ProtectionConcernsAttribute = PartyAttributeDefinition{
		ID:        "ae56b1fd-21f6-480a-9184-091a7093d8b8",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "protectionConcerns",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Beneficiary presents protection concerns"},
						{"es", "La/el beneficiaria/o presenta preocupaciones de protección"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	PhysicalImpairmentAttribute = PartyAttributeDefinition{
		ID:        "cb51b2e8-27da-4375-b85f-c5c107f5d2b4",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "physicalImpairment",
			Type: form.Checkbox,
			Label: i18n.Strings{
				{"en", "Physical impairment"},
				{"es", "Discapacidad física"},
			},
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Would you say you experience some form of physical impairment?"},
						{"es", "¿Diría que experimenta algún tipo de discapacidad física?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	PhysicalImpairmentIntensityAttribute = PartyAttributeDefinition{
		ID:        "98def70b-ee72-40eb-aed1-5a834bf8f579",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "physicalImpairmentIntensity",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Physical impairment intensity"},
				{"es", "Intensidad de la discapacidad física"},
			},
			Description: i18n.Strings{
				{"en", "How would you define the intensity of the physical impairment?"},
				{"es", "¿Cómo definiría la intensidad de la discapacidad física?"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Moderate"},
					{"es", "Moderada"},
				},
				{
					{"en", "Severe"},
					{"es", "Severa"},
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	SensoryImpairmentAttribute = PartyAttributeDefinition{
		ID:        "972c0d7f-8fa9-436d-95ab-6773070bc451",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "sensoryImpairment",
			Type: form.Checkbox,
			Label: i18n.Strings{
				{"en", "Sensory impairment"},
				{"es", "Discapacidad sensorial"},
			},
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Would you say you experience some form of physical impairment?"},
						{"es", "¿Diría que experimenta algún tipo de discapacidad sensorial?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	SensoryImpairmentIntensityAttribute = PartyAttributeDefinition{
		ID:        "b1e6cfac-a8b9-4a0d-a5c7-f164fde99bcc",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "sensoryImpairmentIntensity",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Sensory impairment intensity"},
				{"es", "Intensidad de la discapacidad sensorial"},
			},
			Description: i18n.Strings{
				{"en", "How would you define the intensity of the sensory impairment?"},
				{"es", "¿Cómo definiría la intensidad de la discapacidad sensorial?"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Moderate"},
					{"es", "Moderada"},
				},
				{
					{"en", "Severe"},
					{"es", "Severa"},
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	MentalImpairmentAttribute = PartyAttributeDefinition{
		ID:        "41b7eb87-6488-47e3-a4b0-1422c039d0c7",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "mentalImpairment",
			Type: form.Checkbox,
			Label: i18n.Strings{
				{"en", "Mental impairment"},
				{"es", "Discapacidad mental"},
			},
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Would you say you experience some form of mental impairment?"},
						{"es", "¿Diría que experimenta algún tipo de discapacidad mental?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
	MentalImpairmentIntensityAttribute = PartyAttributeDefinition{
		ID:        "9983188b-4f43-4cd5-a972-fde3a08f4810",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "mentalImpairmentIntensity",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Mental impairment intensity"},
				{"es", "Intensidad de la discapacidad mental"},
			},
			Description: i18n.Strings{
				{"en", "How would you define the intensity of the mental impairment?"},
				{"es", "¿Cómo definiría la intensidad de la discapacidad mental?"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Moderate"},
					{"es", "Moderada"},
				},
				{
					{"en", "Severe"},
					{"es", "Severa"},
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
)

// Global Individual Contact Attributes ---------------------------------------

var (
	PrimaryPhoneNumberAttribute = PartyAttributeDefinition{
		ID:        "8eae83a8-cbc7-4ab2-a21f-d57cb3bb29ff",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "primaryPhoneNumber",
			Label: i18n.Strings{
				{"en", "Primary Phone Number"},
				{"es", "Número de teléfono primario"},
			},
			Type: form.Phone,
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	SecondaryPhoneNumberAttribute = PartyAttributeDefinition{
		ID:        "1f3016af-ab39-422a-beb8-904b68a1619e",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "secondaryPhoneNumber",
			Label: i18n.Strings{
				{"en", "Secondary Phone Number"},
				{"es", "Número de teléfono secundario"},
			},
			Type: form.Phone,
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	TertiaryPhoneNumberAttribute = PartyAttributeDefinition{
		ID:        "4a0ba072-66a5-403f-bea1-35e9427659fb",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "tertiaryPhoneNumber",
			Label: i18n.Strings{
				{"en", "Tertiary Phone Number"},
				{"es", "Número de teléfono terciario"},
			},
			Type: form.Phone,
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}
)

// Uganda Individual Attributes -----------------------------------------------

var (
	UGIdentificationDateAttribute = PartyAttributeDefinition{
		ID:        "c84b8b93-b974-4bec-b9f7-d437446b24a7",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name: "ugIdentificationDate",
			Type: form.Date,
			Label: i18n.Strings{
				{"en", "Date of identification"},
			},
			Description: i18n.Strings{
				{"en", "Date of first interaction with NRC"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGIdentificationLocationAttribute = PartyAttributeDefinition{
		ID:        "06680252-1a1f-4c9d-85dd-56feef20019d",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name: "ugIdentificationLocation",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Location of Identification"},
				{"es", "Ubicación de la identificación"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Kabusu Access Center"},
				},
				{
					{"en", "Nsambya AccessCenter"},
				},
				{
					{"en", "Kisenyi ICLA Center"},
				},
				{
					{"en", "Lukuli ICLA Center"},
				},
				{
					{"en", "Kawempe ICLA Center"},
				},
				{
					{"en", "Ndejje ICLA Center"},
				},
				{
					{"en", "Mengo Field Office"},
				},
				{
					{"en", "Community (Specify location)"},
				},
				{
					{"en", "Home Visit"},
				},
				{
					{"en", "Phone"},
				},
				{
					{"en", "Other (Specify)"},
				},
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

	UGIdentificationSourceAttribute = PartyAttributeDefinition{
		ID:        "a131a0fb-0270-4feb-8fc9-46e7dd6b5acb",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:  "ugIdentificationSource",
			Type:  form.Dropdown,
			Label: i18n.Strings{{"en", "Source of Identification"}},
			Options: []i18n.Strings{
				{{"en", "Walk-in Center"}},
				{{"en", "FFRM Referral"}},
				{{"en", "Internal Referral (Other – Specify)"}},
				{{"en", "ICLA Outreach Team"}},
				{{"en", "External Referral (Community Leader/Contact)"}},
				{{"en", "External Referral (INGO/LNGO)"}},
				{{"en", "External Referral (Other – Specify)"}},
				{{"en", "Self (Telephone)"}},
				{{"en", "Self (Email)"}},
				{{"en", "Internal Referral (Other NRC Sector – Specify)"}},
				{{"en", "CBP Outreach Team"}},
				{{"en", "Other NRC Outreach Team (Specify)"}},
				{{"en", "External Referral (UN Agency)"}},
				{{"en", "External Referral (Government)"}},
				{{"en", "Other – Specify"}}},

			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGAdmin2Attribute = PartyAttributeDefinition{
		ID:        "44dffbc4-7536-42b9-af84-32ea4e9ed493",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:  "ugAdmin2",
			Type:  form.Dropdown,
			Label: i18n.Strings{{"en", "District / Admin 2"}},
			Options: []i18n.Strings{
				{{"en", "ABIM"}},
				{{"en", "ADJUMANI"}},
				{{"en", "ALEBTONG"}},
				{{"en", "AMOLATAR"}},
				{{"en", "AMUDAT"}},
				{{"en", "AMURIA"}},
				{{"en", "AMURU"}},
				{{"en", "APAC"}},
				{{"en", "BUDAKA"}},
				{{"en", "BUGIRI"}},
				{{"en", "BUIKWE"}},
				{{"en", "BUKOMANSIMBI"}},
				{{"en", "BUKWO"}},
				{{"en", "BULAMBULI"}},
				{{"en", "BULIISA"}},
				{{"en", "BUNDIBUGYO"}},
				{{"en", "BUSHENYI"}},
				{{"en", "BUYENDE"}},
				{{"en", "DOKOLO"}},
				{{"en", "BUTAMBALA"}},
				{{"en", "HOIMA"}},
				{{"en", "IGANGA"}},
				{{"en", "KAABONG"}},
				{{"en", "KABALE"}},
				{{"en", "KABAROLE"}},
				{{"en", "KALANGALA"}},
				{{"en", "KALIRO"}},
				{{"en", "KALUNGU"}},
				{{"en", "KAMULI"}},
				{{"en", "KANUNGU"}},
				{{"en", "KAPCHORWA"}},
				{{"en", "KATAKWI"}},
				{{"en", "KAYUNGA"}},
				{{"en", "SHEEMA"}},
				{{"en", "KITGUM"}},
				{{"en", "KOBOKO"}},
				{{"en", "KOLE"}},
				{{"en", "KOTIDO"}},
				{{"en", "KISORO"}},
				{{"en", "KWEEN"}},
				{{"en", "LAMWO"}},
				{{"en", "LIRA"}},
				{{"en", "LUUKA"}},
				{{"en", "LYANTONDE"}},
				{{"en", "MANAFWA"}},
				{{"en", "MASAKA"}},
				{{"en", "MASINDI"}},
				{{"en", "MAYUGE"}},
				{{"en", "MBALE"}},
				{{"en", "MBARARA"}},
				{{"en", "MOROTO"}},
				{{"en", "MOYO"}},
				{{"en", "NAKAPIRIPIRIT"}},
				{{"en", "NAKASEKE"}},
				{{"en", "NAKASONGOLA"}},
				{{"en", "NAMUTUMBA"}},
				{{"en", "NAPAK"}},
				{{"en", "NEBBI"}},
				{{"en", "NGORA"}},
				{{"en", "BUHWEJU"}},
				{{"en", "NTOROKO"}},
				{{"en", "MARACHA"}},
				{{"en", "OTUKE"}},
				{{"en", "OYAM"}},
				{{"en", "PADER"}},
				{{"en", "RUBIRIZI"}},
				{{"en", "SIRONKO"}},
				{{"en", "SOROTI"}},
				{{"en", "WAKISO"}},
				{{"en", "YUMBE"}},
				{{"en", "ZOMBO"}},
				{{"en", "ISINGIRO"}},
				{{"en", "MITOOMA"}},
				{{"en", "KYEGEGWA"}},
				{{"en", "NTUNGAMO"}},
				{{"en", "RUKUNGIRI"}},
				{{"en", "KAMWENGE"}},
				{{"en", "IBANDA"}},
				{{"en", "KASESE"}},
				{{"en", "KIRUHURA"}},
				{{"en", "KYENJOJO"}},
				{{"en", "MUBENDE"}},
				{{"en", "GOMBA"}},
				{{"en", "KIBOGA"}},
				{{"en", "MPIGI"}},
				{{"en", "KYANKWANZI"}},
				{{"en", "KAKUMIRO"}},
				{{"en", "NWOYA"}},
				{{"en", "KIRYANDONGO"}},
				{{"en", "SERERE"}},
				{{"en", "OMORO"}},
				{{"en", "ARUA"}},
				{{"en", "LWENGO"}},
				{{"en", "SEMBABULE"}},
				{{"en", "RAKAI"}},
				{{"en", "MITYANA"}},
				{{"en", "LUWERO"}},
				{{"en", "MUKONO"}},
				{{"en", "KAMPALA"}},
				{{"en", "BUVUMA"}},
				{{"en", "JINJA"}},
				{{"en", "NAMAYINGO"}},
				{{"en", "BUSIA"}},
				{{"en", "BUDUDA"}},
				{{"en", "TORORO"}},
				{{"en", "BUTALEJA"}},
				{{"en", "BUKEDEA"}},
				{{"en", "KUMI"}},
				{{"en", "PALLISA"}},
				{{"en", "KIBUKU"}},
				{{"en", "KABERAMAIDO"}},
				{{"en", "AGAGO"}},
				{{"en", "KAGADI"}},
				{{"en", "KIBAALE"}},
				{{"en", "GULU"}},
				{{"en", "RUBANDA"}},
			},

			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGAdmin3Attribute = PartyAttributeDefinition{
		ID:        "a17ffa5e-5d62-44cd-b89f-438eeba128ac",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:       "ugAdmin3",
			Type:       form.Text,
			Label:      i18n.Strings{{"en", "Subcounty / Admin 3"}},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGAdmin4Attribute = PartyAttributeDefinition{
		ID:        "f867c62a-dcd0-4778-9f4e-7309d044e905",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:       "ugAdmin4",
			Type:       form.Text,
			Label:      i18n.Strings{{"en", "Parish / Admin 4"}},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGAdmin5Attribute = PartyAttributeDefinition{
		ID:        "f0b34ffc-3e15-4195-8e90-a3e1e4b3940c",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:       "ugAdmin5",
			Type:       form.Text,
			Label:      i18n.Strings{{"en", "Village / Admin 5"}},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGNationalityAttribute = PartyAttributeDefinition{
		ID:        "76aab836-73a6-4a1e-9c17-04b8a4c25d8d",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:  "ugNationality",
			Type:  form.Dropdown,
			Label: i18n.Strings{{"en", "Nationality(ies)"}},
			Options: []i18n.Strings{
				{{"en", "Uganda"}},
				{{"en", "Kenya"}},
				{{"en", "Tanzania"}},
				{{"en", "Rwanda"}},
				{{"en", "Burundi"}},
				{{"en", "Democratic Republic of Congo"}},
				{{"en", "South Sudan"}},
				{{"en", "Sudan"}},
				{{"en", "Somalia"}},
				{{"en", "Ethiopia"}},
			},
			Multiple: true,
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGSpokenLanguagesAttribute = PartyAttributeDefinition{
		ID:        "d041cba5-9486-4390-bc2b-ec7fb03d67ff",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:        "ugSpokenLanguages",
			Label:       i18n.Strings{{"en", "Spoken languages"}},
			Description: i18n.Strings{{"en", "What languages does the beneficiary speak?"}},
			Type:        form.Text,
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGPreferredLanguageAttribute = PartyAttributeDefinition{
		ID:        "da27a6e8-abe3-48d5-bfd9-46033e476a09",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:        "ugPreferredLanguage",
			Type:        form.Text,
			Label:       i18n.Strings{{"en", "Preferred Language"}},
			Description: i18n.Strings{{"en", "What language does the beneficiary prefer for communication?"}},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGPhysicalAddressAttribute = PartyAttributeDefinition{
		ID:        "ac2795e8-15a5-42a0-b11f-b9269ff2a309",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:  "ugPhysicalAddress",
			Type:  form.Textarea,
			Label: i18n.Strings{{"en", "Physical address"}},
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGInstructionOnMakingContactAttribute = PartyAttributeDefinition{
		ID:        "4d399cb3-6653-4a61-92eb-331f07e6c395",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name:  "ugInstructionOnMakingContact",
			Type:  form.Textarea,
			Label: i18n.Strings{{"en", "Instructions on contacting the beneficiary"}},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGCanInitiateContactAttribute = PartyAttributeDefinition{
		ID:        "7476fef0-d116-4b94-b981-ac647e16203d",
		CountryID: GlobalCountry.ID,
		FormControl: form.Control{
			Name: "ugCanInitiateContact",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{{"en", "NRC can initiate contact with Beneficiary."}},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGPreferredMeansOfContactAttribute = PartyAttributeDefinition{
		ID:        "1e7f2db9-eb63-46ae-b6d5-5c171a9e2534",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name:  "ugPreferredMeansOfContact",
			Type:  form.Dropdown,
			Label: i18n.Strings{{"en", "Preferred means of contact"}},
			Options: []i18n.Strings{
				{{"en", "Phone Call"}},
				{{"en", "Text message"}},
				{{"en", "WhatsApp"}},
				{{"en", "Signal"}},
				{{"en", "Telegram"}},
				{{"en", "Email"}},
				{{"en", "Home visit"}}},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	UGRequireAnInterpreterAttribute = PartyAttributeDefinition{
		ID:        "9b6ae87d-8935-49aa-9e32-26e7445d1afc",
		CountryID: UgandaCountry.ID,
		FormControl: form.Control{
			Name: "ugRequireAnInterpreter",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{{"en", "This beneficiary requires an interpreter."}},
					Value: "yes",
				},
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
	COPrimaryNationalityAttribute = PartyAttributeDefinition{
		ID:        "d1ee17c5-a7c5-486f-a1e9-be4ec6d65700",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coPrimaryNationality",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Primary nationality"},
				{"es", "Nacionalidad primaria"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Columbia"},
					{"es", "Columbia"},
				},
				{
					{"en", "Venezuela"},
					{"es", "Venezuela"},
				},
				{
					{"en", "Ecuador"},
					{"es", "Ecuador"},
				},
				{
					{"en", "Panama"},
					{"es", "Panamá"},
				},
				{
					{"en", "Costa Rica"},
					{"es", "Costa Rica"},
				},
				{
					{"en", "Honduras"},
					{"es", "Honduras"},
				},
			},
			Multiple: false,
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COSecondaryNationalityAttribute = PartyAttributeDefinition{
		ID:        "74f39024-a318-4c6a-bb07-dfe55679f78f",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coSecondaryNationality",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Secondary nationality"},
				{"es", "Nacionalidad secundaria"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Columbia"},
					{"es", "Columbia"},
				},
				{
					{"en", "Venezuela"},
					{"es", "Venezuela"},
				},
				{
					{"en", "Ecuador"},
					{"es", "Ecuador"},
				},
				{
					{"en", "Panama"},
					{"es", "Panamá"},
				},
				{
					{"en", "Costa Rica"},
					{"es", "Costa Rica"},
				},
				{
					{"en", "Honduras"},
					{"es", "Honduras"},
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COMaritalStatusAttribute = PartyAttributeDefinition{
		ID:        "8bf6b645-20c1-403b-93bc-c05bbc22f570",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coMaritalStatus",
			Label: i18n.Strings{
				{"en", "Marital Status"},
				{"es", "Estado civil"},
			},
			Type: form.Dropdown,
			Options: []i18n.Strings{
				{
					{"en", "Married"},
					{"es", "Casada/o"},
				},
				{
					{"en", "Single"},
					{"es", "Soltera/o"},
				},
				{
					{"en", "Divorced"},
					{"es", "Divorciada/o"},
				},
				{
					{"en", "Separated"},
					{"es", "Separada/o"},
				},
				{
					{"en", "Widowed"},
					{"es", "Viuda/o"},
				},
			},
			Multiple: true,
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COBeneficiaryTypeAttribute = PartyAttributeDefinition{
		ID:        "796e4eb0-56a7-46bb-b81a-9727e674f1f8",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coBeneficiaryType",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Beneficiary type"},
				{"es", "Tipo de beneficiaria/o"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Student"},
					{"es", "Alumna/o"},
				},
				{
					{"en", "Teacher"},
					{"es", "Profesora/o"},
				},
				{
					{"en", "Community leader"},
					{"es", "Líder comunitaria/o"},
				},
				{
					{"en", "Civil servant"},
					{"es", "Funcionaria/o"},
				}},
			Multiple: true,
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COEthnicityAttribute = PartyAttributeDefinition{
		ID:        "fe26bc55-30b7-4c30-97f1-99e90a3367a8",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coEthnicity",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Ethnicity"},
				{"es", "Etnia"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	CORegistrationDateAttribute = PartyAttributeDefinition{
		ID:        "7623b9f3-c29e-479f-872f-bd008a37aca4",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coRegistrationDate",
			Type: form.Date,
			Label: i18n.Strings{
				{"en", "Registration date"},
				{"es", "Fecha de registro"},
			},
			Description: i18n.Strings{
				{"en", "Date of registration with NRC"},
				{"es", "Fecha de registro con el NRC"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	CORegistrationLocationAttribute = PartyAttributeDefinition{
		ID:        "f5ea04e0-7073-45b3-aa9a-a08afaf503da",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coRegistrationLocation",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Location of Registration"},
				{"es", "Lugar de registro"},
			},
			Options: []i18n.Strings{
				{
					{"en", ""},
					{"es", "Viento Libre"},
				},
				{
					{"en", "Other (Specify)"},
					{"es", "Otro (especificar)"},
				},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COSourceOfIdentificationAttribute = PartyAttributeDefinition{
		ID:        "533dd5a3-8ab1-4eb0-9e20-8f4c7b02b2e9",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coSourceOfIdentification",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Source of Identification"},
				{"es", "Fuente de identificación"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Route"},
					{"es", "Calle"},
				},
				{
					{"en", "Shelter"},
					{"es", "Abrigo"},
				},
				{
					{"en", "Protective Space"},
					{"es", "Espacio protector"},
				},
				{
					{"en", "Home"},
					{"es", "Hogar"},
				},
				{
					{"en", "Community"},
					{"es", "Comunidad"},
				}},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COTypeOfSettlementAttribute = PartyAttributeDefinition{
		ID:        "ac56561b-64e4-4d96-bbe8-813a0ed7060c",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coTypeOfSettlement",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Type of settlement"},
				{"es", "Tipo de asentamiento"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COEmergencyCareAttribute = PartyAttributeDefinition{
		ID:        "c425da4b-5af1-4dff-abab-058b1cf9b122",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coEmergencyCare",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{Label: i18n.Strings{
					{"en", "Beneficiary requires emergency care"},
					{"es", "La/el beneficiaria/o requiere atención de emergencia"},
				},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	CODurableSolutionsAttribute = PartyAttributeDefinition{
		ID:        "68241403-dd90-4e26-8d30-70db03b92c95",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coDurableSolutions",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Response is a durable solution?"},
						{"es", "¿La respuesta es una solución duradera?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COHardToReachAttribute = PartyAttributeDefinition{
		ID:        "0c327266-47fb-4557-b2fc-a6e394432254",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coHardToReach",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Is the beneficiary in a hard to reach location?"},
						{"es", "¿Está la/el beneficiaria/o en un lugar de difícil acceso?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAttendedCovid19Attribute = PartyAttributeDefinition{
		ID:        "d59241dc-384b-430d-a8f4-f7851ff28615",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coAttendedCovid19",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Did the beneficiary take part in Covid19 emergency training?"},
						{"es", "¿La/el beneficiaria/o participó en la capacitación de emergencia de Covid19?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COIntroSourceAttribute = PartyAttributeDefinition{
		ID:        "dc7f97a3-b927-438e-9bdd-4374ae09b63a",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coIntroSource",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "How was the beneficiary introduced to NRC?"},
				{"es", "¿Cómo se presentó la/el beneficiaria/o a NRC?"},
			},
		},
		IsPersonallyIdentifiableInfo: true,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin1Attribute = PartyAttributeDefinition{
		ID:        "88a5c89a-9f09-4513-a8cb-f81190f9cc0c",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coAdmin1",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Country / Admin 1"},
				{"es", "País / Admin 1"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Columbia"},
					{"es", "Columbia"},
				},
				{
					{"en", "Venezuela"},
					{"es", "Venezuela"},
				},
				{
					{"en", "Ecuador"},
					{"es", "Ecuador"},
				},
				{
					{"en", "Panama"},
					{"es", "Panamá"},
				},
				{
					{"en", "Costa Rica"},
					{"es", "Costa Rica"},
				},
				{
					{"en", "Honduras"},
					{"es", "Honduras"},
				},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin2Attribute = PartyAttributeDefinition{
		ID:        "491d0ca0-0b63-4860-8e38-8139fcdccf51",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coAdmin2",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "District / Admin 2"},
				{"es", "Departamento / Admin 2"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin3Attribute = PartyAttributeDefinition{
		ID:        "8e69cfdf-935e-43cb-81a0-79ebdda742ec",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coAdmin3",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Subcounty / Admin 3"},
				{"es", "Municipio / Admin 3"},
			},

			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin4Attribute = PartyAttributeDefinition{
		ID:        "cb132ade-f379-42a8-88b0-6c08b375e086",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coAdmin4",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Parish / Admin 4"},
				{"es", "Comuna o Corregimiento / Admin 4"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAdmin5Attribute = PartyAttributeDefinition{
		ID:        "faf65cc6-f5eb-4d18-91ca-00bbd3a3ab8e",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coAdmin5",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Village / Admin 5"},
				{"es", "Barrio o Vereda / Admin 5"},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COJobOrEnterpriseAttribute = PartyAttributeDefinition{
		ID:        "dda85258-9ce2-41e3-a7f8-21b837d65a25",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coJobOrEnterprise",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Do you have a job or enterprise?"},
						{"es", "Ustedes tiene empleo o emprendimiento?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COTypeOfEnterpriseAttribute = PartyAttributeDefinition{
		ID:        "94a9b0d8-8eb7-4165-ae3f-fcf7279da537",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coTypeOfEnterprise",
			Type: form.Dropdown,
			Label: i18n.Strings{
				{"en", "Type of enterprise"},
				{"es", "Tipo de emprendimiento"},
			},
			Options: []i18n.Strings{
				{
					{"en", "Commerce"},
					{"es", "Comercio"},
				},
				{
					{"en", "Production"},
					{"es", "Producción"},
				},
				{
					{"en", "Service"},
					{"es", "Servicio"},
				},
				{
					{"en", "Agriculture"},
					{"es", "Agricultura"},
				},
			},
			Multiple: true,
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COTimeInBusinessAttribute = PartyAttributeDefinition{
		ID:        "31e6f25b-d0a8-47c5-8161-0fdfdb39d430",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coTimeInBusiness",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Length of time in business"},
				{"es", "Tiempo del emprendimiento"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COTypeOfEmploymentAttribute = PartyAttributeDefinition{
		ID:        "fd992e38-2ac1-41ed-9efb-1962174a6438",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coTypeOfEmployment",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Type of employment (type of contract)"},
				{"es", "Tipo de empleo (modalidad de contrato)"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COFormsOfIncomeGenerationAttribute = PartyAttributeDefinition{
		ID:        "0ac65773-d8ef-4f63-907b-9761e1630be8",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coFormsOfIncomeGeneration",
			Type: form.Textarea,
			Label: i18n.Strings{
				{"en", "Forms of income generation in the family"},
				{"es", "Formas de generación de ingresos en la familia"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COLegalRepresentativeNameAttribue = PartyAttributeDefinition{
		ID:        "48a22db2-e97d-44db-869b-0a192697c781",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coLegalRepresentativeName",
			Type: form.Text,
			Label: i18n.Strings{
				{"en", "Name and surname of the legal representative"},
				{"es", "Nombre y apellido del representante jurídico"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COLegalRepresentativeAdditionalInfoAttribute = PartyAttributeDefinition{
		ID:        "f82b59f4-3610-4869-91a3-13308361d153",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coLegalRepresentativeAdditionalInfo",
			Type: form.Textarea,
			Label: i18n.Strings{
				{"en", "Additional information about the legal representative"},
				{"es", "Información adicional sobre el representante jurídico"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COReasonsForRepresentationAttribute = PartyAttributeDefinition{
		ID:        "42673d92-48fb-4426-b13f-104d2625a8ff",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coReasonsForRepresentation",
			Type: form.Textarea,
			Label: i18n.Strings{
				{"en", "Reasons for representation"},
				{"es", "Razones para representar"},
			},
			Validation: form.ControlValidation{Required: true},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COGuardianshipIsLegalAttribute = PartyAttributeDefinition{
		ID:        "5af889e2-c3c9-4ceb-a6a4-d6fff7aa5747",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coGuardianshipIsLegal",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Is the guardianship legal according to national legislation?"},
						{"es", "¿La tutela es legal según la legislación nacional?"},
					},
					Value: "yes",
				},
			},
		},
		IsPersonallyIdentifiableInfo: false,
		PartyTypeIDs: []string{
			IndividualPartyType.ID,
		},
	}

	COAbleToGiveLegalConsentAttribute = PartyAttributeDefinition{
		ID:        "6d96c1e3-9a3c-40c8-93ae-36636ced0b1a",
		CountryID: ColombiaCountry.ID,
		FormControl: form.Control{
			Name: "coAbleToGiveLegalConsent",
			Type: form.Checkbox,
			CheckboxOptions: []form.CheckboxOption{
				{
					Label: i18n.Strings{
						{"en", "Is the person able to give legal consent?"},
						{"es", "¿La persona puede dar su consentimiento de forma legal?"},
					},
					Value: "yes",
				},
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

var TeamNameAttribute = PartyAttributeDefinition{
	ID:        "18f410a3-6fde-45ce-80c7-fc5d92b85870",
	CountryID: GlobalCountry.ID,
	FormControl: form.Control{
		Name: "teamName",
	},
	PartyTypeIDs:                 []string{TeamPartyType.ID},
	IsPersonallyIdentifiableInfo: false,
}

var CountryNameAttribute = PartyAttributeDefinition{
	ID:        "e011d638-864b-496e-b3e5-af89d0278e1e",
	CountryID: GlobalCountry.ID,
	FormControl: form.Control{
		Name: "countryName",
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
