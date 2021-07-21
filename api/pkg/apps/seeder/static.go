package seeder

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"strings"
)

func caseType(id, name, partyTypeID, teamID string, template *cms.CaseTemplate) cms.CaseType {
	ct := cms.CaseType{
		ID:          id,
		Name:        name,
		PartyTypeID: partyTypeID,
		TeamID:      teamID,
		Template:    template,
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

func kase(id, caseTypeID, createdByID, partyID, teamID string, done bool, formData *cms.CaseTemplate) cms.Case {
	k := cms.Case{
		ID:         id,
		CaseTypeID: caseTypeID,
		CreatorID:  createdByID,
		PartyID:    partyID,
		TeamID:     teamID,
		Done:       done,
		FormData:   formData,
	}
	cases = append(cases, k)
	return k
}

var (
	teams         []iam.Team
	individuals   []iam.Individual
	staffers      []iam.Staff
	memberships   []iam.Membership
	relationships []iam.Relationship
	caseTypes     []cms.CaseType
	cases         []cms.Case

	// Teams
	UgandaProtectionTeam    = team("ac9b8d7d-d04d-4850-9a7f-3f93324c0d1e", "Uganda Protection Team")
	UgandaICLATeam          = team("a43f84d5-3f8a-48c4-a896-5fb0fcd3e42b", "Uganda ICLA Team")
	UgandaCoreAdminTeam     = team("814fc372-08a6-4e6b-809b-30ebb51cb268", "Uganda Core Admin Team")
	MozambiqueEducationTeam = team("80606eb4-b53a-4fda-be12-e9806e11d44a", "Mozambique Education Team")

	// Case Templates for Uganda
	// - Kampala Response Team
	UGSituationAnalysis = &cms.CaseTemplate{
		FormElements: []cms.FormElement{
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life?",
					ID:          "safeDiginifiedLife",
					Description: "Probe for description",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?",
					ID:          "challengesBarriers",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?",
					ID:          "solutions",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "If we were to work together on this, what could we do together? What would make the most difference for you?",
					ID:          "workTogether",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}
	UGIndividualAssessment = &cms.CaseTemplate{
		FormElements: []cms.FormElement{
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Which service has the individual requested as a starting point of support?",
					ID:          "serviceStartingPoint",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "What other services has the individual requested/identified?",
					ID:          "otherServices",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "What is the perceived priority response level of the individual",
					ID:          "perceivedPriority",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}
	UGReferral = &cms.CaseTemplate{
		FormElements: []cms.FormElement{
			{
				Type: cms.TextInput,
				Attributes: cms.FormElementAttribute{
					Label:       "Date of Referral",
					ID:          "dateOfReferral",
					Description: "",
				},
			},
			{
				Type: cms.Dropdown,
				Attributes: cms.FormElementAttribute{
					Label:       "Urgency",
					ID:          "urgency",
					Description: "",
					Options:     []string{"Very Urgent", "Urgent", "Not Urgent"},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Dropdown,
				Attributes: cms.FormElementAttribute{
					Label:       "Type of Referral",
					ID:          "typeOfReferral",
					Description: "",
					Options:     []string{"Internal", "External"},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Services/assistance requested",
					ID:          "servicesRequested",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Reason for referral",
					ID:          "reasonForReferral",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Checkbox,
				Attributes: cms.FormElementAttribute{
					Label:       "Does the beneficiary have any restrictions to be referred?",
					ID:          "referralRestrictions",
					Description: "",
					CheckboxOptions: []cms.CheckboxOption{
						{
							Label: "Has restrictions?",
						},
					},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Dropdown,
				Attributes: cms.FormElementAttribute{
					Label:       "Means of Referral",
					ID:          "meansOfReferral",
					Description: "",
					Options:     []string{"Phone", "E-mail", "Personal meeting", "Other"},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Means and terms of receiving feedback from the client",
					ID:          "meansOfFeedback",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.TextInput,
				Attributes: cms.FormElementAttribute{
					Label:       "Deadline for receiving feedback from the client",
					ID:          "deadlineForFeedback",
					Description: "",
				},
			},
		},
	}
	UGExternalReferralFollowup = &cms.CaseTemplate{
		FormElements: []cms.FormElement{
			{
				Type: cms.Checkbox,
				Attributes: cms.FormElementAttribute{
					Label:       "Was the referral accepted by the other provider?",
					ID:          "referralAccepted",
					Description: "",
					CheckboxOptions: []cms.CheckboxOption{
						{
							Label: "Referral accepted",
						},
					},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Provide any pertinent details on service needs / requests.",
					ID:          "pertinentDetails",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}
	// - Kampala ICLA Team
	UGICLAIndividualIntake = &cms.CaseTemplate{
		FormElements: []cms.FormElement{
			{
				Type: cms.Dropdown,
				Attributes: cms.FormElementAttribute{
					Label:       "Modality of service delivery",
					ID:          "modality",
					Description: "",
					Options:     []string{"ICLA Legal Aid Centre", "Mobile visit", "Home visit", "Transit Centre", "Hotline", "Other"},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Dropdown,
				Attributes: cms.FormElementAttribute{
					Label:       "Living situation",
					ID:          "livingSituation",
					Description: "",
					Options:     []string{"Lives alone", "Lives with family", "Hosted by relatives"},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Comment on living situation",
					ID:          "commentLivingSituation",
					Description: "Additional information, observations, concerns, etc.",
					Placeholder: "",
				},
			},
			{
				Type: cms.Dropdown,
				Attributes: cms.FormElementAttribute{
					Label:       "How did you learn about ICLA services?",
					ID:          "iclaServiceDiscovery",
					Description: "",
					Options:     []string{"ICLA in-person information session", "ICLA social media campaign, activities, brochures", "ICLA text messages", "Another beneficiary/friend/relative", "Another organisation", "General social media", "NRC employee", "State authority", "Other"},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Vulnerability data",
					ID:          "vulnerability",
					Description: "As needed within a particular context and required for the case",
					Placeholder: "",
				},
			},
			{
				Type: cms.TextInput,
				Attributes: cms.FormElementAttribute{
					Label:       "Full name of representative",
					ID:          "representativeName",
					Description: "Lawyer or other person",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Other personal information",
					ID:          "otherInformation",
					Description: "Other personal data as needed to identify the representative within the particular context",
					Placeholder: "",
				},
			},
			{
				Type: cms.TextInput,
				Attributes: cms.FormElementAttribute{
					Label:       "Reason for representative",
					ID:          "representativeReason",
					Description: "",
				},
			},
			{
				Type: cms.Checkbox,
				Attributes: cms.FormElementAttribute{
					Label:       "Is the guardianship legal as per national legislation?",
					ID:          "guardianshipIsLegal",
					Description: "If 'yes', attach/upload the legal assessment. If 'no', request or assist in identifying an appropriate legal guardian to represent beneficiary",
					CheckboxOptions: []cms.CheckboxOption{
						{
							Label: "Guardianship is legal",
						},
					},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Checkbox,
				Attributes: cms.FormElementAttribute{
					Label:       "Does the beneficiary have the legal capacity to consent?",
					ID:          "capacityToConsent",
					Description: "",
					CheckboxOptions: []cms.CheckboxOption{
						{
							Label: "Beneficiary has legal capacity to consent",
						},
					},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
		},
	}
	UGICLACaseAssessment = &cms.CaseTemplate{
		FormElements: []cms.FormElement{
			{
				Type: cms.Dropdown,
				Attributes: cms.FormElementAttribute{
					Label:       "Type of service",
					ID:          "serviceType",
					Description: "",
					Options:     []string{"Legal counselling", "Legal assistance"},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.TextInput,
				Attributes: cms.FormElementAttribute{
					Label:       "Thematic area",
					ID:          "thematicArea",
					Description: "Applicable Thematic Area related to the problem",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Fact and details of the problem",
					ID:          "details",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Checkbox,
				Attributes: cms.FormElementAttribute{
					Label:       "Other parties involved",
					ID:          "otherPartiesInvolved",
					Description: "Are there any other parties involved in the case",
					CheckboxOptions: []cms.CheckboxOption{
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
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Checkbox,
				Attributes: cms.FormElementAttribute{
					Label:       "Previous/existing lawyer working on the case",
					ID:          "previousOrExistingLawyer",
					Description: "Does the client have a previous or existing lawyer working on his/her case?",
					CheckboxOptions: []cms.CheckboxOption{
						{
							Label: "Previous lawyer",
						},
						{
							Label: "Existing lawyer",
						},
					},
				},
				Validation: cms.FormElementValidation{
					Required: true,
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Previous or existing lawyer details",
					ID:          "previousOrExistingLawyerDetails",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "What actions have been taken to solve the problem, if any?",
					ID:          "actionsTaken",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Related to this problem, are there any cases pending before a court or administrative body?",
					ID:          "pendingCourtCases",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "If there are cases pending before a court or administrative body, are there any deadlines that need to be met?",
					ID:          "pendingCourtCaseDeadlines",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: cms.Textarea,
				Attributes: cms.FormElementAttribute{
					Label:       "Is there any conflict of interest involved?",
					ID:          "conflictOfInterest",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}

	// Case Types for Uganda
	// - Kampala Response Team
	UGSituationalAnalysisCaseType      = caseType("0ae90b08-6944-48dc-8f30-5cb325292a8c", "Situational Analysis (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGSituationAnalysis)
	UGIndividualAssessmentCaseType     = caseType("2f909038-0ce4-437b-af17-72fc5d668b49", "Individual Assessment (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGIndividualAssessment)
	UGReferralCaseType                 = caseType("ecdaf47f-6fa9-48c8-9d10-6324bf932ed7", "Referral (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGReferral)
	UGExternalReferralFollowupCaseType = caseType("2a1b670c-6336-4364-b89d-0e65fc771659", "External Referral Followup (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGExternalReferralFollowup)
	// - Kampala ICLA Team
	UGICLAIndividualIntakeCaseType = caseType("31fb6d03-2374-4bea-9374-48fc10500f81", "ICLA Individual Intake (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLAIndividualIntake)
	UGICLACaseAssessmentCaseType   = caseType("bbf820de-8d10-49eb-b8c9-728993ab0b73", "ICLA Case Assessment (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLACaseAssessment)

	// Individuals
	JohnDoe     = individual("c529d679-3bb6-4a20-8f06-c096f4d9adc1", "John", "Doe", "12/02/1978", "Refugee", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "No", "Yes", "Moderate", "No", "", "No", "", "Kenya", "Kiswahili, English", "English", "123 Main Street, Kampala", "0123456789", "", "Email", "No")
	MaryPoppins = individual("bbf539fd-ebaa-4438-ae4f-8aca8b327f42", "Mary", "Poppins", "12/02/1978", "Internally Displaced Person", "Female", "Yes", "https://link-to-consent.proof", "No", "No", "No", "No", "", "No", "", "No", "", "Uganda", "Rukiga, English", "Rukiga", "901 First Avenue, Kampala", "0123456789", "", "Telegram", "Yes")
	BoDiddley   = individual("26335292-c839-48b6-8ad5-81271ee51e7b", "Bo", "Diddley", "12/02/1978", "Host Community", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "Yes", "No", "", "No", "", "No", "", "Somalia", "Somali, Arabic, English", "English", "101 Main Street, Kampala", "0123456789", "", "Whatsapp", "No")

	// Individuals (Staff)
	Stephen  = staff(individual("066a0268-fdc6-495a-9e4b-d60cfae2d81a", "Stephen", "Kabagambe", "12/02/1978", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Colette  = staff(individual("93f9461f-31da-402e-8988-6e0100ecaa24", "Colette", "le Jeune", "12/02/1978", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Courtney = staff(individual("14c014d9-f433-4508-b33d-dc45bf86690b", "Courtney", "Lare", "12/02/1978", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))

	// Memberships
	StevenMembership   = membership("862690ee-87f0-4f95-aa1e-8f8a2f2fd54a", Stephen, UgandaCoreAdminTeam)
	ColetteMembership  = membership("9d4abef9-0be0-4750-81ab-0524a412c049", Colette, UgandaProtectionTeam)
	CourtneyMembership = membership("83c5e73a-5947-4d7e-996c-14a2a7b1c850", Courtney, MozambiqueEducationTeam)

	// Cases
	BoDiddleySituationAnalysis = kase("dba43642-8093-4685-a197-f8848d4cbaaa", UGSituationalAnalysisCaseType.ID, Colette.ID, BoDiddley.ID, UgandaProtectionTeam.ID, false, &cms.CaseTemplate{
		FormElements: []cms.FormElement{
			{
				Type: "textarea",
				Attributes: cms.FormElementAttribute{
					Label:       "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life?",
					ID:          "safeDiginifiedLife",
					Description: "Probe for description",
					Value: []string{
						"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life.",
					},
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.FormElementAttribute{
					Label: "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?",
					ID:    "challengesBarriers",
					Value: []string{
						"Some of the barriers I face are communication gaps between myself and refugee tenants. We are attempting to deal with these challenges by using google translate.",
					},
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.FormElementAttribute{
					Label: "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?",
					ID:    "solutions",
					Value: []string{
						"A qualified interpreter, who knows the legal context could help us to agree on contractual matters.",
					},
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.FormElementAttribute{
					Label: "If we were to work together on this, what could we do together? What would make the most difference for you?",
					ID:    "workTogether",
					Value: []string{
						"NRC could provide a translator and a legal representative to ease contract negotiations",
					},
					Description: "",
					Placeholder: "",
				},
			},
		},
	})
)
