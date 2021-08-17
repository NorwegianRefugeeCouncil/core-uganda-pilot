package seeder

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/registrationctrl"
	"strings"
)

func caseType(id, name, partyTypeID, teamID string, template *cms.CaseTemplate, intakeCaseType bool) cms.CaseType {
	ct := cms.CaseType{
		ID:             id,
		Name:           name,
		PartyTypeID:    partyTypeID,
		TeamID:         teamID,
		Template:       template,
		IntakeCaseType: intakeCaseType,
	}
	caseTypes = append(caseTypes, ct)
	return ct
}

func team(id, name string) iam.Team {
	t := iam.Team{
		ID:   id,
		Name: name,
	}
	teams = append(teams, t)
	return t
}

func country(id, name string) iam.Country {
	t := iam.Country{
		ID:   id,
		Name: name,
	}
	countries = append(countries, t)
	return t
}

func individual(
	id string,
	firstName string,
	lastName string,
	birthDate string,
	displacementStatus string,
	gender string,
	consent string,
	consentProof string,
	anonymous string,
	minor string,
	protectionConcerns string,
	physicalImpairment string,
	physicalImpairmentIntensity string,
	sensoryImpairment string,
	sensoryImpairmentIntensity string,
	mentalImpairment string,
	mentalImpairmentIntensity string,
	nationality string,
	spokenLanguages string,
	preferredLanguage string,
	physicalAddress string,
	primaryPhoneNumber string,
	secondaryPhoneNumber string,
	preferredMeansOfContact string,
	requireAnInterpreter string,
) iam.Individual {
	var i = iam.Individual{
		Party: &iam.Party{
			ID: id,
			PartyTypeIDs: []string{
				iam.IndividualPartyType.ID,
			},
			Attributes: map[string][]string{
				iam.FirstNameAttribute.ID:                   {firstName},
				iam.LastNameAttribute.ID:                    {lastName},
				iam.EMailAttribute.ID:                       {strings.ToLower(firstName) + "." + strings.ToLower(lastName) + "@email.com"},
				iam.BirthDateAttribute.ID:                   {birthDate},
				iam.DisplacementStatusAttribute.ID:          {displacementStatus},
				iam.GenderAttribute.ID:                      {gender},
				iam.ConsentToNrcDataUseAttribute.ID:         {consent},
				iam.ConsentToNrcDataUseProofAttribute.ID:    {consentProof},
				iam.AnonymousAttribute.ID:                   {anonymous},
				iam.MinorAttribute.ID:                       {minor},
				iam.ProtectionConcernsAttribute.ID:          {protectionConcerns},
				iam.PhysicalImpairmentAttribute.ID:          {physicalImpairment},
				iam.PhysicalImpairmentIntensityAttribute.ID: {physicalImpairmentIntensity},
				iam.SensoryImpairmentAttribute.ID:           {sensoryImpairment},
				iam.SensoryImpairmentIntensityAttribute.ID:  {sensoryImpairmentIntensity},
				iam.MentalImpairmentAttribute.ID:            {mentalImpairment},
				iam.MentalImpairmentIntensityAttribute.ID:   {mentalImpairmentIntensity},
				iam.NationalityAttribute.ID:                 {nationality},
				iam.SpokenLanguagesAttribute.ID:             {spokenLanguages},
				iam.PreferredLanguageAttribute.ID:           {preferredLanguage},
				iam.PhysicalAddressAttribute.ID:             {physicalAddress},
				iam.PrimaryPhoneNumberAttribute.ID:          {primaryPhoneNumber},
				iam.SecondaryPhoneNumberAttribute.ID:        {secondaryPhoneNumber},
				iam.PreferredMeansOfContactAttribute.ID:     {preferredMeansOfContact},
				iam.RequireAnInterpreterAttribute.ID:        {requireAnInterpreter},
			},
		},
	}
	individuals = append(individuals, i)
	return i
}

func staff(individual iam.Individual) iam.Individual {
	individual.AddPartyType(iam.StaffPartyType.ID)
	return individual
}

func membership(id string, individual iam.Individual, team iam.Team) iam.Membership {
	m := iam.Membership{
		ID:           id,
		TeamID:       team.ID,
		IndividualID: individual.ID,
	}
	memberships = append(memberships, m)
	return m
}

func nationality(id string, team iam.Team, country iam.Country) iam.Nationality {
	m := iam.Nationality{
		ID:        id,
		CountryID: country.ID,
		TeamID:    team.ID,
	}
	nationalities = append(nationalities, m)
	return m
}

func kase(id, caseTypeID, createdByID, partyID, teamID string, done bool, form *cms.CaseTemplate, intakeCase bool) cms.Case {

	k := cms.Case{
		ID:         id,
		CaseTypeID: caseTypeID,
		CreatorID:  createdByID,
		PartyID:    partyID,
		TeamID:     teamID,
		Done:       done,
		Template:   form,
		IntakeCase: intakeCase,
	}
	cases = append(cases, k)
	return k
}

var (
	teams         []iam.Team
	individuals   []iam.Individual
	staffers      []iam.Staff
	memberships   []iam.Membership
	countries     []iam.Country
	nationalities []iam.Nationality
	relationships []iam.Relationship
	caseTypes     []cms.CaseType
	cases         []cms.Case

	// Countries
	Germany = country("02680685-806e-4386-b6a5-95c4af1fc141", "Germany")
	Uganda  = country("062a7fe9-b9cc-4fbc-837e-138a15242007", "Uganda")

	// Teams
	UgandaProtectionTeam = team("ac9b8d7d-d04d-4850-9a7f-3f93324c0d1e", "Uganda Protection Team")
	UgandaICLATeam       = team("a43f84d5-3f8a-48c4-a896-5fb0fcd3e42b", "Uganda ICLA Team")
	UgandaCoreAdminTeam  = team("814fc372-08a6-4e6b-809b-30ebb51cb268", "Uganda Core Admin Team")

	// Case Templates for Uganda
	// - Kampala Response Team
	UGSituationAnalysis = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life?",
					Name:        "safeDiginifiedLife",
					Description: "Probe for description",
					Placeholder: "",
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?",
					Name:        "challengesBarriers",
					Description: "",
					Placeholder: "",
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?",
					Name:        "solutions",
					Description: "",
					Placeholder: "",
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "If we were to work together on this, what could we do together? What would make the most difference for you?",
					Name:        "workTogether",
					Description: "",
					Placeholder: "",
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
		},
	}
	UGIndividualResponse = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.TaxonomyInput,
				Attributes: form.FormElementAttributes{
					Label:       "Which service has the individual requested as a starting point of support?",
					Name:        "serviceStartingPoint",
					Description: "Add the taxonomies of the services requested as a starting point one by one, by selecting the relevant options from the dropdowns below.",
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.TaxonomyInput,
				Attributes: form.FormElementAttributes{
					Label:       "What other services has the individual requested/identified?",
					Name:        "otherServices",
					Description: "Add the taxonomies of the other services requested one by one, by selecting the relevant options from the dropdowns below.",
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Text,
				Attributes: form.FormElementAttributes{
					Label: "What is the perceived priority response level of the individual",
					Name:  "perceivedPriority",
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
		},
	}
	UGReferral = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.Text,
				Attributes: form.FormElementAttributes{
					Label:       "Date of Referral",
					Name:        "dateOfReferral",
					Description: "",
				},
			},
			{
				Type: form.Dropdown,
				Attributes: form.FormElementAttributes{
					Label:       "Urgency",
					Name:        "urgency",
					Description: "",
					Options:     []string{"Very Urgent", "Urgent", "Not Urgent"},
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Dropdown,
				Attributes: form.FormElementAttributes{
					Label:       "Type of Referral",
					Name:        "typeOfReferral",
					Description: "",
					Options:     []string{"Internal", "External"},
				},
				Validation: form.FormElementValidation{
					Required: false,
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Services/assistance requested",
					Name:        "servicesRequested",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Reason for referral",
					Name:        "reasonForReferral",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Checkbox,
				Attributes: form.FormElementAttributes{
					Label:       "Does the beneficiary have any restrictions to be referred?",
					Name:        "referralRestrictions",
					Description: "",
					CheckboxOptions: []form.CheckboxOption{
						{
							Label: "Has restrictions?",
						},
					},
				},
			},
			{
				Type: form.Dropdown,
				Attributes: form.FormElementAttributes{
					Label:       "Means of Referral",
					Name:        "meansOfReferral",
					Description: "",
					Options:     []string{"Phone", "E-mail", "Personal meeting", "Other"},
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Means and terms of receiving feedback from the client",
					Name:        "meansOfFeedback",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Text,
				Attributes: form.FormElementAttributes{
					Label:       "Deadline for receiving feedback from the client",
					Name:        "deadlineForFeedback",
					Description: "",
				},
			},
		},
	}
	UGExternalReferralFollowup = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.Checkbox,
				Attributes: form.FormElementAttributes{
					Label:       "Was the referral accepted by the other provider?",
					Name:        "referralAccepted",
					Description: "",
					CheckboxOptions: []form.CheckboxOption{
						{
							Label: "Referral accepted",
						},
					},
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Provide any pertinent details on service needs / requests.",
					Name:        "pertinentDetails",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}
	// - Kampala ICLA Team
	UGICLAIndividualIntake = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.Dropdown,
				Attributes: form.FormElementAttributes{
					Label:       "Modality of service delivery",
					Name:        "modality",
					Description: "",
					Options:     []string{"ICLA Legal Aid Centre", "Mobile visit", "Home visit", "Transit Centre", "Hotline", "Other"},
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Dropdown,
				Attributes: form.FormElementAttributes{
					Label:       "Living situation",
					Name:        "livingSituation",
					Description: "",
					Options:     []string{"Lives alone", "Lives with family", "Hosted by relatives"},
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Comment on living situation",
					Name:        "commentLivingSituation",
					Description: "Additional information, observations, concerns, etc.",
					Placeholder: "",
				},
			},
			{
				Type: form.Dropdown,
				Attributes: form.FormElementAttributes{
					Label:       "How did you learn about ICLA services?",
					Name:        "iclaServiceDiscovery",
					Description: "",
					Options:     []string{"ICLA in-person information session", "ICLA social media campaign, activities, brochures", "ICLA text messages", "Another beneficiary/friend/relative", "Another organisation", "General social media", "NRC employee", "State authority", "Other"},
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Vulnerability data",
					Name:        "vulnerability",
					Description: "As needed within a particular context and required for the case",
					Placeholder: "",
				},
			},
			{
				Type: form.Text,
				Attributes: form.FormElementAttributes{
					Label:       "Full name of representative",
					Name:        "representativeName",
					Description: "Lawyer or other person",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Other personal information",
					Name:        "otherInformation",
					Description: "Other personal data as needed to identify the representative within the particular context",
					Placeholder: "",
				},
			},
			{
				Type: form.Text,
				Attributes: form.FormElementAttributes{
					Label:       "Reason for representative",
					Name:        "representativeReason",
					Description: "",
				},
			},
			{
				Type: form.Checkbox,
				Attributes: form.FormElementAttributes{
					Label:       "Is the guardianship legal as per national legislation?",
					Name:        "guardianshipIsLegal",
					Description: "If 'yes', attach/upload the legal assessment. If 'no', request or assist in identifying an appropriate legal guardian to represent beneficiary",
					CheckboxOptions: []form.CheckboxOption{
						{
							Label: "Guardianship is legal",
						},
					},
				},
			},
			{
				Type: form.Checkbox,
				Attributes: form.FormElementAttributes{
					Label:       "Does the beneficiary have the legal capacity to consent?",
					Name:        "capacityToConsent",
					Description: "",
					CheckboxOptions: []form.CheckboxOption{
						{
							Label: "Beneficiary has legal capacity to consent",
						},
					},
				},
			},
		},
	}
	UGICLACaseAssessment = &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.Dropdown,
				Attributes: form.FormElementAttributes{
					Label:       "Type of service",
					Name:        "serviceType",
					Description: "",
					Options:     []string{"Legal counselling", "Legal assistance"},
				},
				Validation: form.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: form.Text,
				Attributes: form.FormElementAttributes{
					Label:       "Thematic area",
					Name:        "thematicArea",
					Description: "Applicable Thematic Area related to the problem",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Fact and details of the problem",
					Name:        "details",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Checkbox,
				Attributes: form.FormElementAttributes{
					Label:       "Other parties involved",
					Name:        "otherPartiesInvolved",
					Description: "Are there any other parties involved in the case",
					CheckboxOptions: []form.CheckboxOption{
						{
							Label: "Landlord",
						},
						{
							Label: "Lawyer",
						},
						{
							Label: "Relative",
						},
						{
							Label: "Other",
						},
					},
				},
			},
			{
				Type: form.Checkbox,
				Attributes: form.FormElementAttributes{
					Label:       "Previous/existing lawyer working on the case",
					Name:        "previousOrExistingLawyer",
					Description: "Does the client have a previous or existing lawyer working on his/her case?",
					CheckboxOptions: []form.CheckboxOption{
						{
							Label: "Previous lawyer",
						},
						{
							Label: "Existing lawyer",
						},
					},
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Previous or existing lawyer details",
					Name:        "previousOrExistingLawyerDetails",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "What actions have been taken to solve the problem, if any?",
					Name:        "actionsTaken",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Related to this problem, are there any cases pending before a court or administrative body?",
					Name:        "pendingCourtCases",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "If there are cases pending before a court or administrative body, are there any deadlines that need to be met?",
					Name:        "pendingCourtCaseDeadlines",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Is there any conflict of interest involved?",
					Name:        "conflictOfInterest",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}

	// - Kampala Response Team
	UGSituationalAnalysisCaseType      = caseType("0ae90b08-6944-48dc-8f30-5cb325292a8c", "Situational Analysis (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGSituationAnalysis, true)
	UGIndividualResponseCaseType       = caseType("2f909038-0ce4-437b-af17-72fc5d668b49", "Response (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGIndividualResponse, true)
	UGReferralCaseType                 = caseType("ecdaf47f-6fa9-48c8-9d10-6324bf932ed7", "Referral (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGReferral, false)
	UGExternalReferralFollowupCaseType = caseType("2a1b670c-6336-4364-b89d-0e65fc771659", "External Referral Followup (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGExternalReferralFollowup, false)
	// - Kampala ICLA Team
	UGICLAIndividualIntakeCaseType = caseType("31fb6d03-2374-4bea-9374-48fc10500f81", "ICLA Individual Intake (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLAIndividualIntake, true)
	UGICLACaseAssessmentCaseType   = caseType("bbf820de-8d10-49eb-b8c9-728993ab0b73", "ICLA Case Assessment (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLACaseAssessment, false)

	UgandaRegistrationFlow = registrationctrl.RegistrationFlow{
		// TODO Country
		TeamID: "",
		Steps: []registrationctrl.Step{{
			Type:  registrationctrl.IndividualAttributes,
			Label: "New individual intake",
			Ref:   "",
		}, {
			Type:  registrationctrl.CaseType,
			Label: "Situation Analysis",
			Ref:   UGSituationalAnalysisCaseType.ID,
		}, {
			Type:  registrationctrl.CaseType,
			Label: "Individual Response",
			Ref:   UGIndividualResponseCaseType.ID,
		}},
	}

	// Individuals
	JohnDoe     = individual("c529d679-3bb6-4a20-8f06-c096f4d9adc1", "John", "Doe", "1983-04-23", "Refugee", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "No", "Yes", "Moderate", "No", "", "No", "", "Kenya", "Kiswahili, English", "English", "123 Main Street, Kampala", "0123456789", "", "Email", "No")
	MaryPoppins = individual("bbf539fd-ebaa-4438-ae4f-8aca8b327f42", "Mary", "Poppins", "1983-04-23", "Internally Displaced Person", "Female", "Yes", "https://link-to-consent.proof", "No", "No", "No", "No", "", "No", "", "No", "", "Uganda", "Rukiga, English", "Rukiga", "901 First Avenue, Kampala", "0123456789", "", "Telegram", "Yes")
	BoDiddley   = individual("26335292-c839-48b6-8ad5-81271ee51e7b", "Bo", "Diddley", "1983-04-23", "Host Community", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "Yes", "No", "", "No", "", "No", "", "Somalia", "Somali, Arabic, English", "English", "101 Main Street, Kampala", "0123456789", "", "Whatsapp", "No")

	// Individuals (Staff)
	Stephen  = staff(individual("066a0268-fdc6-495a-9e4b-d60cfae2d81a", "Stephen", "Kabagambe", "1983-04-23", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Colette  = staff(individual("93f9461f-31da-402e-8988-6e0100ecaa24", "Colette", "le Jeune", "1983-04-23", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Courtney = staff(individual("14c014d9-f433-4508-b33d-dc45bf86690b", "Courtney", "Lare", "1983-04-23", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))

	// Memberships
	StevenMembership   = membership("862690ee-87f0-4f95-aa1e-8f8a2f2fd54a", Stephen, UgandaCoreAdminTeam)
	ColetteMembership  = membership("9d4abef9-0be0-4750-81ab-0524a412c049", Colette, UgandaProtectionTeam)
	CourtneyMembership = membership("83c5e73a-5947-4d7e-996c-14a2a7b1c850", Courtney, UgandaProtectionTeam)

	// Nationalities
	DTeamNationality                = nationality("9c1c1f2d-67f5-41cc-a752-534f031c05f9", DTeam, Germany)
	UgandaCoreAdminTeamNationality  = nationality("0987460d-c906-43cd-b7fd-5e7afca0d93e", UgandaCoreAdminTeam, Uganda)
	UgandaProtectionTeamNationality = nationality("b58e4d26-fe8e-4442-8449-7ec4ca3d9066", UgandaProtectionTeam, Uganda)
	UgandaICLATeamNationality       = nationality("23e3eb5e-592e-42e2-8bbf-ee097d93034c", UgandaICLATeam, Uganda)

	// Cases
	BoDiddleySituationAnalysis = kase("dba43642-8093-4685-a197-f8848d4cbaaa", UGSituationalAnalysisCaseType.ID, Colette.ID, BoDiddley.ID, UgandaProtectionTeam.ID, true, &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life?",
					Name:        "safeDiginifiedLife",
					Description: "Probe for description",
					Value: []string{
						"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life.",
					},
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?",
					Name:  "challengesBarriers",
					Value: []string{
						"Some of the barriers I face are communication gaps between myself and refugee tenants. We are attempting to deal with these challenges by using google translate.",
					},
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?",
					Name:  "solutions",
					Value: []string{
						"A qualified interpreter, who knows the legal context could help us to agree on contractual matters.",
					},
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "If we were to work together on this, what could we do together? What would make the most difference for you?",
					Name:  "workTogether",
					Value: []string{
						"NRC could provide a translator and a legal representative to ease contract negotiations",
					},
					Description: "",
					Placeholder: "",
				},
			},
		},
	},
		true)

	BoDiddleyIndividualAssessment = kase("3ea8c121-bdf0-46a0-86a8-698dc4abc872", UGIndividualResponseCaseType.ID, Colette.ID, BoDiddley.ID, UgandaProtectionTeam.ID, true, &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.TaxonomyInput,
				Attributes: form.FormElementAttributes{
					Label: "Which service has the individual requested as a starting point of support?",
					Name:  "serviceStartingPoint",
					Value: []string{
						"ICLA - Counselling - Legal Identity",
					},
					Description: "Add the taxonomies of the services requested as a starting point one by one, by selecting the relevant options from the dropdowns below.",
					Placeholder: "",
				},
			},
			{
				Type: form.TaxonomyInput,
				Attributes: form.FormElementAttributes{
					Label: "What other services has the individual requested/identified?",
					Name:  "otherServices",
					Value: []string{
						"Protection - Individual Targeted Protection - Individual Protection Assistance",
					},
					Description: "Add the taxonomies of the other services requested one by one, by selecting the relevant options from the dropdowns below.",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "What is the perceived priority response level of the individual",
					Name:        "perceivedPriority",
					Description: "",
					Placeholder: "",
					Value: []string{
						"High",
					},
				},
			},
		},
	}, true)

	MaryPoppinsSituationAnalysis = kase("4f7708ed-240a-423f-9bd1-839542e65833", UGSituationalAnalysisCaseType.ID, Colette.ID, MaryPoppins.ID, UgandaProtectionTeam.ID, true, &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life?",
					Name:        "safeDiginifiedLife",
					Description: "Probe for description",
					Value: []string{
						"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life.",
					},
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?",
					Name:  "challengesBarriers",
					Value: []string{
						"Some of the barriers I face are communication gaps between myself and refugee tenants. We are attempting to deal with these challenges by using google translate.",
					},
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?",
					Name:  "solutions",
					Value: []string{
						"A qualified interpreter, who knows the legal context could help us to agree on contractual matters.",
					},
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "If we were to work together on this, what could we do together? What would make the most difference for you?",
					Name:  "workTogether",
					Value: []string{
						"NRC could provide a translator and a legal representative to ease contract negotiations",
					},
					Description: "",
					Placeholder: "",
				},
			},
		},
	},
		true)

	MaryPoppinsIndividualAssessment = kase("45b4a637-c610-4ab9-afe6-4e958c36a96f", UGIndividualResponseCaseType.ID, Colette.ID, MaryPoppins.ID, UgandaProtectionTeam.ID, true, &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.TaxonomyInput,
				Attributes: form.FormElementAttributes{
					Label: "Which service has the individual requested as a starting point of support?",
					Name:  "serviceStartingPoint",
					Value: []string{
						"ICLA - Counselling - Legal Identity",
					},
					Description: "Add the taxonomies of the services requested as a starting point one by one, by selecting the relevant options from the dropdowns below.",
					Placeholder: "",
				},
			},
			{
				Type: form.TaxonomyInput,
				Attributes: form.FormElementAttributes{
					Label: "What other services has the individual requested/identified?",
					Name:  "otherServices",
					Value: []string{
						"Protection - Individual Targeted Protection - Individual Protection Assistance",
					},
					Description: "Add the taxonomies of the other services requested one by one, by selecting the relevant options from the dropdowns below.",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "What is the perceived priority response level of the individual",
					Name:        "perceivedPriority",
					Description: "",
					Value: []string{
						"High",
					},
					Placeholder: "",
				},
			},
		},
	}, true)

	JohnDoesSituationAnalysis = kase("43140381-8166-4fb3-9ac5-339082920ade", UGSituationalAnalysisCaseType.ID, Colette.ID, JohnDoe.ID, UgandaProtectionTeam.ID, true, &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label:       "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life?",
					Name:        "safeDiginifiedLife",
					Description: "Probe for description",
					Value: []string{
						"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life.",
					},
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?",
					Name:  "challengesBarriers",
					Value: []string{
						"Some of the barriers I face are communication gaps between myself and refugee tenants. We are attempting to deal with these challenges by using google translate.",
					},
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?",
					Name:  "solutions",
					Value: []string{
						"A qualified interpreter, who knows the legal context could help us to agree on contractual matters.",
					},
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "If we were to work together on this, what could we do together? What would make the most difference for you?",
					Name:  "workTogether",
					Value: []string{
						"NRC could provide a translator and a legal representative to ease contract negotiations",
					},
					Description: "",
					Placeholder: "",
				},
			},
		},
	},
		true)

	JohnDoeIndividualAssessment = kase("65e02e79-1676-4745-9890-582e3d67d13f", UGIndividualResponseCaseType.ID, Colette.ID, JohnDoe.ID, UgandaProtectionTeam.ID, true, &cms.CaseTemplate{
		FormElements: []form.FormElement{
			{
				Type: form.TaxonomyInput,
				Attributes: form.FormElementAttributes{
					Label: "Which service has the individual requested as a starting point of support?",
					Name:  "serviceStartingPoint",
					Value: []string{
						"ICLA - Counselling - Legal Identity",
					},
					Description: "Add the taxonomies of the services requested as a starting point one by one, by selecting the relevant options from the dropdowns below.",
					Placeholder: "",
				},
			},
			{
				Type: form.TaxonomyInput,
				Attributes: form.FormElementAttributes{
					Label: "What other services has the individual requested/identified?",
					Name:  "otherServices",
					Value: []string{
						"Protection - Individual Targeted Protection - Individual Protection Assistance",
					},
					Description: "Add the taxonomies of the other services requested one by one, by selecting the relevant options from the dropdowns below.",
					Placeholder: "",
				},
			},
			{
				Type: form.Textarea,
				Attributes: form.FormElementAttributes{
					Label: "What is the perceived priority response level of the individual",
					Name:  "perceivedPriority",
					Value: []string{
						"High",
					},
					Description: "",
					Placeholder: "",
				},
			},
		},
	}, true)
)
