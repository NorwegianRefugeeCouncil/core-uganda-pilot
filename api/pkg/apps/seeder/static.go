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
				iam.FirstNameAttribute.ID: {firstName},
				iam.LastNameAttribute.ID:  {lastName},
				iam.EMailAttribute.ID:     {strings.ToLower(firstName) + "." + strings.ToLower(lastName) + "@email.com"},
				iam.BirthDateAttribute.ID: {birthDate},
				iam.DisplacementStatusAttribute.ID: {displacementStatus},
				iam.GenderAttribute.ID: {gender},
				iam.ConsentToNrcDataUseAttribute.ID: {consent},
				iam.ConsentToNrcDataUseProofAttribute.ID: {consentProof},
				iam.AnonymousAttribute.ID: {anonymous},
				iam.MinorAttribute.ID: {minor},
				iam.ProtectionConcernsAttribute.ID: {protectionConcerns},
				iam.PhysicalImpairmentAttribute.ID: {physicalImpairment},
				iam.PhysicalImpairmentIntensityAttribute.ID: {physicalImpairmentIntensity},
				iam.SensoryImpairmentAttribute.ID: {sensoryImpairment},
				iam.SensoryImpairmentIntensityAttribute.ID: {sensoryImpairmentIntensity},
				iam.MentalImpairmentAttribute.ID: {mentalImpairment},
				iam.MentalImpairmentIntensityAttribute.ID: {mentalImpairmentIntensity},
				iam.NationalityAttribute.ID: {nationality},
				iam.SpokenLanguagesAttribute.ID: {spokenLanguages},
				iam.PreferredLanguageAttribute.ID: {preferredLanguage},
				iam.PhysicalAddressAttribute.ID: {physicalAddress},
				iam.PrimaryPhoneNumberAttribute.ID: {primaryPhoneNumber},
				iam.SecondaryPhoneNumberAttribute.ID: {secondaryPhoneNumber},
				iam.PreferredMeansOfContactAttribute.ID: {preferredMeansOfContact},
				iam.RequireAnInterpreterAttribute.ID: {requireAnInterpreter},
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

func kase(id, caseTypeID, createdByID, partyID, teamID string, done bool) cms.Case {
	k := cms.Case{
		ID:         id,
		CaseTypeID: caseTypeID,
		CreatorID:  createdByID,
		PartyID:    partyID,
		TeamID:     teamID,
		Done:       done,
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
	KampalaResponseTeam = team("ac9b8d7d-d04d-4850-9a7f-3f93324c0d1e", "Kampala Response Team")
	KampalaICLATeam     = team("a43f84d5-3f8a-48c4-a896-5fb0fcd3e42b", "Kampala ICLA Team")
	NairobiResponseTeam = team("814fc372-08a6-4e6b-809b-30ebb51cb268", "Nairobi Response Team")
	NairobiICLATeam     = team("80606eb4-b53a-4fda-be12-e9806e11d44a", "Nairobi ICLA Team")

	// Case Templates
	Legal = &cms.CaseTemplate{
		FormElements: []cms.CaseTemplateFormElement{
			{
				Type: "dropdown",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Legal satus",
					ID:          "legalStatus",
					Description: "What is the beneficiary's current legal status?",
					Options:     []string{"Citizen", "Permanent resident", "Accepted refugee", "Asylum seeker", "Undetermined"},
				},
				Validation: cms.CaseTemplateFormElementValidation{
					Required: true,
				},
			},
			{
				Type: "checkbox",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Qualified services",
					ID:          "qualifiedServices",
					Description: "What services does the beneficiary qualify for?",
					CheckboxOptions: []cms.CaseTemplateCheckboxOption{
						{
							Label: "Counselling",
						},
						{
							Label: "Representation",
						},
						{
							Label: "Arbitration",
						},
					},
				},
				Validation: cms.CaseTemplateFormElementValidation{
					Required: true,
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Notes",
					ID:          "notes",
					Description: "Additional information, observations, concerns, etc.",
					Placeholder: "Type here",
				},
			},
			{
				Type: "textinput",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Project number",
					ID:          "projectNumber",
					Description: "Enter the beneficiaries project number, if any",
				},
			},
		},
	}

	// Case Templates for Uganda
	// - Kampala Response Team
	UGSituationAnalysis = &cms.CaseTemplate{
		FormElements: []cms.CaseTemplateFormElement{
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life?",
					ID:          "safeDiginifiedLife",
					Description: "Probe for description",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?",
					ID:          "challengesBarriers",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?",
					ID:          "solutions",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "If we were to work together on this, what could we do together? What would make the most difference for you?",
					ID:          "workTogether",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}
	UGIndividualAssessment = &cms.CaseTemplate{
		FormElements: []cms.CaseTemplateFormElement{
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Which service has the individual requested as a starting point of support?",
					ID:          "serviceStartingPoint",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "What other services has the individual requested/identified?",
					ID:          "otherServices",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "What is the perceived priority response level of the individual",
					ID:          "perceivedPriority",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}
	UGReferral = &cms.CaseTemplate{
		FormElements: []cms.CaseTemplateFormElement{
			{
				Type: "textinput",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Date of Referral",
					ID:          "dateOfReferral",
					Description: "",
				},
			},
			{
				Type: "dropdown",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Urgency",
					ID:          "urgency",
					Description: "",
					Options:     []string{"Very Urgent", "Urgent", "Not Urgent"},
				},
				Validation: cms.CaseTemplateFormElementValidation{
					Required: true,
				},
			},
			{
				Type: "dropdown",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Type of Referral",
					ID:          "typeOfReferral",
					Description: "",
					Options:     []string{"Internal", "External"},
				},
				Validation: cms.CaseTemplateFormElementValidation{
					Required: true,
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Services/assistance requested",
					ID:          "servicesRequested",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Reason for referral",
					ID:          "reasonForReferral",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "checkbox",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Does the beneficiary have any restrictions to be referred?",
					ID:          "referralRestrictions",
					Description: "",
					CheckboxOptions: []cms.CaseTemplateCheckboxOption{
						{
							Label: "Has restrictions?",
						},
					},
				},
				Validation: cms.CaseTemplateFormElementValidation{
					Required: true,
				},
			},
			{
				Type: "dropdown",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Means of Referral",
					ID:          "meansOfReferral",
					Description: "",
					Options:     []string{"Phone", "E-mail", "Personal meeting", "Other"},
				},
				Validation: cms.CaseTemplateFormElementValidation{
					Required: true,
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Means and terms of receiving feedback from the client",
					ID:          "meansOfFeedback",
					Description: "",
					Placeholder: "",
				},
			},
			{
				Type: "textinput",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Deadline for receiving feedback from the client",
					ID:          "deadlineForFeedback",
					Description: "",
				},
			},
		},
	}
	UGExternalReferralFollowup = &cms.CaseTemplate{
		FormElements: []cms.CaseTemplateFormElement{
			{
				Type: "checkbox",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Was the referral accepted by the other provider?",
					ID:          "referralAccepted",
					Description: "",
					CheckboxOptions: []cms.CaseTemplateCheckboxOption{
						{
							Label: "Referral accepted",
						},
					},
				},
				Validation: cms.CaseTemplateFormElementValidation{
					Required: true,
				},
			},
			{
				Type: "textarea",
				Attributes: cms.CaseTemplateFormElementAttribute{
					Label:       "Provide any pertinent details on service needs / requests.",
					ID:          "pertinentDetails",
					Description: "",
					Placeholder: "",
				},
			},
		},
	}
	// - Kampala ICLA Team

	// Case Types for Uganda
	// - Kampala Response Team
	UGSituationalAnalysisCaseType = caseType("0ae90b08-6944-48dc-8f30-5cb325292a8c", "Situational Analysis (Kampala Response)", iam.IndividualPartyType.ID, KampalaResponseTeam.ID, UGSituationAnalysis)
	UGIndividualAssessmentCaseType = caseType("2f909038-0ce4-437b-af17-72fc5d668b49", "Individual Assessment (Kampala Response)", iam.IndividualPartyType.ID, KampalaResponseTeam.ID, UGIndividualAssessment)
	UGReferralCaseType = caseType("ecdaf47f-6fa9-48c8-9d10-6324bf932ed7", "Referral (Kampala Response)", iam.IndividualPartyType.ID, KampalaResponseTeam.ID, UGReferral)
	UGExternalReferralFollowupCaseType = caseType("2a1b670c-6336-4364-b89d-0e65fc771659", "External Referral Followup (Kampala Response)", iam.IndividualPartyType.ID, KampalaResponseTeam.ID, UGExternalReferralFollowup)
	// - Kampala ICLA Team

	// Individuals
	JohnDoe     = individual("c529d679-3bb6-4a20-8f06-c096f4d9adc1", "John", "Doe", "12/02/1978", "Refugee", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "No", "Yes", "Moderate", "No", "", "No", "", "Kenya", "Kiswahili, English", "English", "123 Main Street, Kampala", "0123456789", "", "Email", "No")
	MaryPoppins = individual("bbf539fd-ebaa-4438-ae4f-8aca8b327f42", "Mary", "Poppins", "12/02/1978", "Internally Displaced Person", "Female", "Yes", "https://link-to-consent.proof", "No", "No", "No", "No", "", "No", "", "No", "", "Uganda", "Rukiga, English", "Rukiga", "901 First Avenue, Kampala", "0123456789", "", "Telegram", "Yes")
	BoDiddley   = individual("26335292-c839-48b6-8ad5-81271ee51e7b", "Bo", "Diddley", "12/02/1978", "Host Community", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "Yes", "No", "", "No", "", "No", "", "Somalia", "Somali, Arabic, English", "English", "101 Main Street, Kampala", "0123456789", "", "Whatsapp", "No")
	Howell      = staff(individual("066a0268-fdc6-495a-9e4b-d60cfae2d81a", "Howell", "Jorg", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Birdie      = staff(individual("ac9015ac-686f-4719-9c3d-bf3d1cae00ea", "Birdie", "Tifawt", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Charis      = staff(individual("ce7ae69c-9f6a-413b-96bf-5808d0da92cd", "Charis", "Timothy", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Danyial     = staff(individual("8da05c97-12c2-4b43-b022-dc79be7dc3a0", "Danyial", "Hrodebert", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Devi        = staff(individual("6414895a-ce60-4647-b491-baeb54a76f26", "Devi", "Malvina", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Levan       = staff(individual("5d4c4302-ad8e-45ab-bd4c-e1ac25ae972e", "Levan", "Elija", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Lisbeth     = staff(individual("a5d4dab0-90d3-474d-afe6-46d04ca3caba", "Lisbeth", "Furkan", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Liadan      = staff(individual("c9ce906d-87ba-4123-bb74-7a73664e6778", "Liadan", "Jordaan", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Muhammad    = staff(individual("818206ea-0b5e-4ed9-b47e-db31566d10c0", "Muhammad", "Annemarie", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Dardanos    = staff(individual("7921756a-8759-4589-8a83-ad98f8aa22c7", "Dardanos", "Rilla", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Jana        = staff(individual("c7ca3a4d-0e96-4e5c-8c32-6750d0312706", "Jana", "Nurul", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Simeon      = staff(individual("78663ffb-dbaa-4362-83b6-7319d6469caa", "Simeon", "Tumelo", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Sayen       = staff(individual("29a20d76-dd37-471f-b9ec-9ab08f61d1ed", "Sayen", "Gezabele", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Veniaminu   = staff(individual("051a46b2-1ef4-4c86-bd2f-9306daedec7e", "Veniaminu", "Ye-Jun", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Loan        = staff(individual("f2a5d586-6865-40ea-a3db-7c729516b32b", "Loan", "Daniel", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Reece       = staff(individual("bdeb7e66-9129-467e-abc0-51ab2df7f222", "Reece", "Hyakinthos", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Svetlana    = staff(individual("afdd8b5c-b9b4-41e1-a015-7e0beb33f10b", "Svetlana", "Cerdic", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Kyleigh     = staff(individual("12d6a293-d923-47c6-9bc1-441934bb79c5", "Kyleigh", "Jayma", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Hermina     = staff(individual("dafee423-49c0-4fbf-b2f9-a42276c0cfce", "Hermina", "Magnus", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Leela       = staff(individual("65410229-ad41-4c17-88f2-13e9a56a0fe8", "Leela", "Cynebald", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Jovan       = staff(individual("bf22e83b-cfef-4c8a-b74e-f0cef6b27147", "Jovan", "Lynette", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Bor         = staff(individual("e350e394-091f-469c-a217-488b27b113a3", "Bor", "Lora", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Aldwin      = staff(individual("fdb6a682-8eb6-4565-879b-835a76384fe0", "Aldwin", "Colin", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))
	Trophimos   = staff(individual("bb800fe3-85a7-4c90-b8f2-cd0354825f56", "Trophimos", "Wiebke", "12/02/1978", "Refugee", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""))

	// Memberships
	HowellMembership    = membership("862690ee-87f0-4f95-aa1e-8f8a2f2fd54a", Howell, KampalaResponseTeam)
	BirdieMembership    = membership("5fa34696-80e2-4a3e-ad26-5cc06535f67c", Birdie, KampalaResponseTeam)
	CharisMembership    = membership("72db9abe-8297-4995-b30d-bafe79f01d90", Charis, KampalaResponseTeam)
	DanyialMembership   = membership("24f6acfd-d1dd-40f8-a5b7-2b9a74d4f70b", Danyial, KampalaResponseTeam)
	DeviMembership      = membership("9811ad8e-febd-4ea0-8dba-0188eec52b94", Devi, KampalaICLATeam)
	LevanMembership     = membership("341e0a25-352e-43cb-9e5a-ffc6ce373c61", Levan, KampalaICLATeam)
	LisbethMembership   = membership("102deee5-5cf9-49c9-a9a0-99b2bde85eae", Lisbeth, KampalaICLATeam)
	BorMembership       = membership("196ad5b6-3375-4acd-83ca-1b4d6f1de19c", Bor, KampalaICLATeam)
	LiadanMembership    = membership("7f6087dd-d4d3-4a92-8c22-90bddc3b28a8", Liadan, NairobiResponseTeam)
	MuhammadMembership  = membership("113a0595-b3e3-422c-8a14-0d60ff71bb17", Muhammad, NairobiResponseTeam)
	DardanosMembership  = membership("fbbd25e9-5a2b-46d3-a8b5-a52fab5801d7", Dardanos, NairobiResponseTeam)
	JanaMembership      = membership("da3795dc-dbd9-4213-bfaa-c10764c664ba", Jana, NairobiResponseTeam)
	SimeonMembership    = membership("a7f5ad21-dd00-4d6e-92be-68c186793935", Simeon, NairobiICLATeam)
	SayenMembership     = membership("9d99551b-5cd3-4948-8695-2ee73c79f13c", Sayen, NairobiICLATeam)
	VeniaminuMembership = membership("ea2b4a53-2968-405c-9c26-8618adba6540", Veniaminu, NairobiICLATeam)
	LoanMembership      = membership("340d8740-b029-41ef-9db6-2bdf991c3ed3", Loan, NairobiICLATeam)

	// Cases
	//DomesticAbuse    = kase("dba43642-8093-4685-a197-f8848d4cbaaa", GenderViolence.ID, Birdie.ID, MaryPoppins.ID, KampalaResponseTeam.ID, false)
	//MonthlyAllowance = kase("47499762-c189-4a74-9156-7969f899073b", FinancialAssistInd.ID, Birdie.ID, JohnDoe.ID, KampalaResponseTeam.ID, false)
	//ChildCare        = kase("8fb5f755-85eb-4d91-97a9-fdf86c01df25", Childcare.ID, Birdie.ID, BoDiddley.ID, KampalaResponseTeam.ID, true)
)
