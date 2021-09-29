package seeder

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/form"
	"github.com/nrc-no/core/pkg/i18n"
	"github.com/nrc-no/core/pkg/registrationctrl"
)

func caseType(id, name, partyTypeID, teamID string, form form.Form, intakeCaseType bool) cms.CaseType {
	ct := cms.CaseType{
		ID:             id,
		Name:           name,
		PartyTypeID:    partyTypeID,
		TeamID:         teamID,
		Form:           form,
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

func individual(id string, fullName string, displayName string, birthDate string, email string, displacementStatus string, gender string, consent string, consentProof string, anonymous string, minor string, protectionConcerns string, physicalImpairment string, physicalImpairmentIntensity string, sensoryImpairment string, sensoryImpairmentIntensity string, mentalImpairment string, mentalImpairmentIntensity string, nationality string, spokenLanguages string, preferredLanguage string, physicalAddress string, primaryPhoneNumber string, secondaryPhoneNumber string, preferredMeansOfContact string, requireAnInterpreter string) iam.Individual {
	var i = iam.Individual{
		Party: &iam.Party{
			ID: id,
			PartyTypeIDs: []string{
				iam.IndividualPartyType.ID,
			},
			Attributes: map[string][]string{
				iam.FullNameAttribute.ID:                    {fullName},
				iam.DisplayNameAttribute.ID:                 {displayName},
				iam.EMailAttribute.ID:                       {email + "@email.com"},
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
				iam.UGNationalityAttribute.ID:               {nationality},
				iam.UGSpokenLanguagesAttribute.ID:           {spokenLanguages},
				iam.UGPreferredLanguageAttribute.ID:         {preferredLanguage},
				iam.UGPhysicalAddressAttribute.ID:           {physicalAddress},
				iam.PrimaryPhoneNumberAttribute.ID:          {primaryPhoneNumber},
				iam.SecondaryPhoneNumberAttribute.ID:        {secondaryPhoneNumber},
				iam.UGPreferredMeansOfContactAttribute.ID:   {preferredMeansOfContact},
				iam.UGRequireAnInterpreterAttribute.ID:      {requireAnInterpreter},
			},
		},
	}
	individuals = append(individuals, i)
	return i
}

func ugandaIndividual(
	individual iam.Individual,
	identificationDate string,
	identificationLocation string,
	identificationSource string,
	admin2 string,
	admin3 string,
	admin4 string,
	admin5 string,
) iam.Individual {
	individual.Attributes.Add(iam.UGIdentificationDateAttribute.ID, identificationDate)
	individual.Attributes.Add(iam.UGIdentificationLocationAttribute.ID, identificationLocation)
	individual.Attributes.Add(iam.UGIdentificationSourceAttribute.ID, identificationSource)
	individual.Attributes.Add(iam.UGAdmin2Attribute.ID, admin2)
	individual.Attributes.Add(iam.UGAdmin3Attribute.ID, admin3)
	individual.Attributes.Add(iam.UGAdmin4Attribute.ID, admin4)
	individual.Attributes.Add(iam.UGAdmin5Attribute.ID, admin5)
	return individual
}

func staff(individual iam.Individual) iam.Individual {
	individual.AddPartyType(iam.StaffPartyType.ID)
	return individual
}

func beneficiary(individual iam.Individual) iam.Individual {
	individual.AddPartyType(iam.BeneficiaryPartyType.ID)
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

func kase(id, createdByID, partyID, teamID string, caseType cms.CaseType, done, intakeCase bool, formData map[string][]string) cms.Case {
	k := cms.Case{
		ID:         id,
		CaseTypeID: caseType.ID,
		CreatorID:  createdByID,
		PartyID:    partyID,
		TeamID:     teamID,
		Done:       done,
		Form:       caseType.Form,
		FormData:   formData,
		IntakeCase: intakeCase,
	}
	cases = append(cases, k)
	return k
}

func identificationDocumentType(id, name string) iam.IdentificationDocumentType {
	idt := iam.IdentificationDocumentType{
		ID:   id,
		Name: name,
	}
	identificationDocumentTypes = append(identificationDocumentTypes, idt)
	return idt
}

func identificationDocument(id, partyId, documentNumber, identificationDocumentTypeId string) iam.IdentificationDocument {
	newId := iam.IdentificationDocument{
		ID:                           id,
		PartyID:                      partyId,
		DocumentNumber:               documentNumber,
		IdentificationDocumentTypeID: identificationDocumentTypeId,
	}
	identificationDocuments = append(identificationDocuments, newId)
	return newId
}

var (
	teams                       []iam.Team
	individuals                 []iam.Individual
	staffers                    []iam.Staff
	memberships                 []iam.Membership
	countries                   []iam.Country
	nationalities               []iam.Nationality
	relationships               []iam.Relationship
	caseTypes                   []cms.CaseType
	cases                       []cms.Case
	identificationDocumentTypes []iam.IdentificationDocumentType
	identificationDocuments     []iam.IdentificationDocument

	// Teams
	KampalaProtectionTeam = team("ac9b8d7d-d04d-4850-9a7f-3f93324c0d1e", "Kampala Protection Team")
	KampalaICLATeam       = team("a43f84d5-3f8a-48c4-a896-5fb0fcd3e42b", "Kampala ICLA Team")
	KampalaCOTeam         = team("814fc372-08a6-4e6b-809b-30ebb51cb268", "Kampala CO Team")
	ColombiaTeam          = team("a6bc6436-fcea-4738-bde8-593e6480e1ad", "Colombia Team")

	// Case Templates for Uganda
	// - Kampala Response Team
	UGSituationAnalysis = form.Form{
		Controls: []form.Control{
			{
				Name:        "safeDignifiedLife",
				Type:        form.Textarea,
				Label:       i18n.Strings{{"en", "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life?"}},
				Description: i18n.Strings{{"en", "Probe for description"}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "challengesBarriers",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?"}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "solutions",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?"}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "workTogether",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "If we were to work together on this, what could we do together? What would make the most difference for you?"}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
		},
	}
	UGIndividualResponse = form.Form{
		Controls: []form.Control{
			{
				Name:        "servicesStartingPoint",
				Type:        form.Taxonomy,
				Label:       i18n.Strings{{"en", "Which service has the individual requested as a starting point of support?"}},
				Description: i18n.Strings{{"en", "Add the taxonomies of the services requested as a starting point one by one, by selecting the relevant options from the dropdowns below."}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:        "commentStartingPoint",
				Type:        form.Textarea,
				Label:       i18n.Strings{{"en", "Comment on service the individual requested as a starting point of support?"}},
				Description: i18n.Strings{{"en", "Additional information, observations, concerns, etc."}},
			},
			{
				Name: "otherServices",
				Type: form.Taxonomy,

				Label:       i18n.Strings{{"en", "What other services has the individual requested/identified?"}},
				Description: i18n.Strings{{"en", "Add the taxonomies of the other services requested one by one, by selecting the relevant options from the dropdowns below."}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:        "commentOtherServices",
				Type:        form.Textarea,
				Label:       i18n.Strings{{"en", "Comment on other services the individual requested/identified?"}},
				Description: i18n.Strings{{"en", "Additional information, observations, concerns, etc."}},
			},
			{
				Name:  "perceivedPriority",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "What is the perceived priority response level of the individual"}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
		},
	}
	UGReferral = form.Form{
		Controls: []form.Control{
			{
				Name:  "dateOfReferral",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Date of Referral"}},
			},
			{
				Name:    "ugency",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Urgency"}},
				Options: []i18n.Strings{{{"en", "Very Urgent"}}, {{"en", "Urgent"}}, {{"en", "Not Urgent"}}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:    "typeOfReferral",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Type of Referral"}},
				Options: []i18n.Strings{{{"en", "Internal"}}, {{"en", "External"}}},
			},
			{
				Name:  "servicesRequested",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Services/assistance requested"}},
			},
			{
				Name:  "readonforReferral",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Reason for referral"}},
			},
			{
				Name:  "referralRestrictions",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Does the beneficiary have any restrictions to be referred?"}},
			},
			{
				Name:    "meansOfReferral",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Means of Referral"}},
				Options: []i18n.Strings{{{"en", "Phone"}}, {{"en", "E-mail"}}, {{"en", "Personal meeting"}}, {{"en", "Other"}}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "meansOfFeedback",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Means and terms of receiving feedback from the client"}},
			},
			{
				Name:  "deadlineForFeedback",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Deadline for receiving feedback from the client"}},
			},
		},
	}
	UGExternalReferralFollowup = form.Form{
		Controls: []form.Control{
			{
				Name:  "referralAccepted",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Was the referral accepted by the other provider?"}},
			},
			{
				Name:  "pertinentDetails",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Provide any pertinent details on service needs / requests."}},
			},
		},
	}

	UGICLACaseAssessment = form.Form{
		Controls: []form.Control{
			{
				Name:    "serviceType",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Type of service"}},
				Options: []i18n.Strings{{{"en", "Legal counselling"}}, {{"en", "Legal assistance"}}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:        "thematicArea",
				Type:        form.Text,
				Label:       i18n.Strings{{"en", "Thematic area"}},
				Description: i18n.Strings{{"en", "Applicable Thematic Area related to the problem"}},
			},
			{
				Name:  "problemDetails",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Fact and details of the problem"}},
			},
			{
				Name: "otherPartiesInvolved",
				Type: form.Checkbox,

				Label:       i18n.Strings{{"en", "Other parties involved"}},
				Description: i18n.Strings{{"en", "Are there any other parties involved in the case"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Landlord"}},
					},
					{
						Label: i18n.Strings{{"en", "Lawyer"}},
					},
					{
						Label: i18n.Strings{{"en", "Relative"}},
					},
					{
						Label: i18n.Strings{{"en", "Other"}},
					},
				},
			},
			{
				Name:  "previousOrExistingLawyer",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Previous/existing lawyer working on the case"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Previous lawyer"}},
					},
					{
						Label: i18n.Strings{{"en", "Existing lawyer"}},
					},
				},
			},
			{
				Name:  "previousOrExistingLawyerDetails",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Previous or existing lawyer details"}},
			},
			{
				Name:  "actionsTaken",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What actions have been taken to solve the problem, if any?"}},
			},
			{
				Name:  "pendingCourtCases",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Related to this problem, are there any cases pending before a court or administrative body?"}},
			},
			{
				Name:  "pendingCourtDeadlines",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "If there are cases pending before a court or administrative body, are there any deadlines that need to be met?"}},
			},
			{
				Name:  "conflictOfInterest",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Is there any conflict of interest involved?"}},
			},
		},
	}

	UGICLAFollowUp = form.Form{
		Controls: []form.Control{
			{
				Name:  "dateOfFollowUp",
				Type:  form.Date,
				Label: i18n.Strings{{"en", "Date of follow up"}},
			},
			{
				Name:     "actionPoints",
				Type:     form.Dropdown,
				Multiple: true,
				Label:    i18n.Strings{{"en", "What action points did you follow up on?"}},
				Options: []i18n.Strings{
					{{"en", "(1) In-person interview / meeting with beneficiary"}},
					{{"en", "(2) Phone conversation with a beneficiary"}},
					{{"en", "(3) Discussion with supervisor/team leader"}},
					{{"en", "(4) Conducting legal analysis, including the study of judicial practice"}},
					{{"en", "(5) Preparing letters, inquiries to various authorities"}},
					{{"en", "(6) Drafting of other legal documents (such leases or contracts)"}},
					{{"en", "(7) Lodging of a court application"}},
					{{"en", "(8) Attending of court session/hearing"}},
					{{"en", "(9) Review of the decision/appeal"}},
					{{"en", "(10) Execution of the court decision"}},
					{{"en", "(11) Negotiation"}},
					{{"en", "(12) Follow up with relevant administrative authority or other entities"}},
					{{"en", "(13) Accompaniment"}},
					{{"en", "(14) Other"}},
				},
			},
			{
				Name:  "notes",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Notes from the follow up undertaken"}},
			},
			{
				Name:  "documentLink",
				Type:  form.URL,
				Label: i18n.Strings{{"en", "Document link"}},
			},
			{
				Name:  "dateOfReceipt",
				Type:  form.Date,
				Label: i18n.Strings{{"en", "Date of receipt"}},
			},
		},
		Sections: []form.Section{
			{
				Title: i18n.Strings{{"en", "FOLLOW UP"}},
				ControlNames: []string{
					"dateOfFollowUp",
					"actionPoints",
					"notes",
				},
			},
			{
				Title: i18n.Strings{{"en", "DOCUMENTS RECEIVED FROM BENEFICIARY"}},
				ControlNames: []string{
					"documentLink",
					"dateOfReceipt",
				},
			},
		},
	}

	UGICLAIntake = form.Form{
		Sections: []form.Section{
			{
				Title: []i18n.LocaleString{{"en", "Legal Issue"}},
				ControlNames: []string{
					"legalIssueDescription",
					"legalIssue",
					"otherLegalIssueDescription",
					"legalActionsTaken",
				},
			},
			{
				Title: []i18n.LocaleString{{"en", "Information of beneficiary's representative"}},
				ControlNames: []string{
					"hasRepresentative",
					"representativeFullName",
					"reasonForRepresentative",
					"isLegalGuardianship",
					"courtOrder",
				},
			},
			{
				Title: i18n.Strings{{"en", "RSD"}},
				ControlNames: []string{
					"individualDisplacementStatus",
					"isAtRiskStateless",
					"statelessRiskDescription",
					"rsdDocuments",
					"rsdComment",
					"rsdIssues",
				},
			},
			{
				Title: i18n.Strings{{"en", "HLP"}},
				ControlNames: []string{
					"specificHLPConcern",
				},
			},
			{
				Title: i18n.Strings{{"en", "Housing"}},
				ControlNames: []string{
					"indStay",
					"hasRentalAgreement",
					"agreementKind",
					"hasEvictionRisk",
					"evictionDoc",
					"evictionComment",
				},
			},
			{
				Title: i18n.Strings{{"en", "Land"}},
				ControlNames: []string{
					"isLegalOwner",
					"natureLandTenancy",
					"natureLandTenure",
					"hasLandEvictionRisk",
					"landEvictionProof",
					"specificLandIssues",
				},
			},
			{
				Title: i18n.Strings{{"en", "Property"}},
				ControlNames: []string{
					"propertyNature",
					"hasLegalOwnershipOfProperty",
					"propertyOwnershipProof",
					"propertyAcquisitionProof",
				},
			},
			{
				Title: i18n.Strings{{"en", "LCD"}},
				ControlNames: []string{
					"documentationChallenges",
					"typeOfDocument",
					"lcdActionsTaken",
				},
			},
			{
				Title: i18n.Strings{{"en", "ELP"}},
				ControlNames: []string{
					"employmentOrBusinessChallenge",
					"employmentChallenges",
					"typeOfEmploymentAgreement",
					"elpActionsTaken",
					"businessRelatedChallenges",
					"businessRegistrationNeeds",
				},
			},
		},
		Controls: []form.Control{
			{
				Name:  "legalIssueDescription",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What/Describe legal/concerns are you facing"}},
			},
			{
				Name:  "legalIssue",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Select the legal issue of concern"}},
				Options: []i18n.Strings{
					{{"en", "RSD"}},
					{{"en", "ELP"}},
					{{"en", "HLP"}},
					{{"en", "IDP registration"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "otherLegalIssueDescription",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "If other, please specify"}},
			},
			{
				Name:  "legalActionsTaken",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What action has been taken to solve the problem, if any"}},
			},
			// subsection Information of beneficiary's representative
			{
				Name:  "hasRepresentative",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Is there a representative (lawyer or another person/Legal Guardian/Other) for this individual?"}},
			},
			// TODO dependant fields
			{
				Name:  "representativeFullName",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Full name of representative"}},
			},
			{
				Name:  "reasonForRepresentative",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Reason for representative (instead of beneficiary)"}},
			},
			{
				Name:  "isLegalGuardianship",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Is the guardianship legal as per national legislation?"}},
			},
			{
				Name:  "courtOrder",
				Type:  form.File,
				Label: []i18n.LocaleString{{"en", "Attach/upload the legal/court order"}},
			},
			//  subsection "RSD"
			{
				Name:  "individualDisplacementStatus",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "What is the individual's displacement status?"}},
				Options: []i18n.Strings{
					{{"en", "Unregistered asylum seeker"}},
					{{"en", "Registered asylum seeker"}},
					{{"en", "Refugee"}},
				},
			},
			{
				Name:  "isAtRiskStateless",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Are you at risk of being stateless?"}},
			},
			{
				Name:  "statelessRiskDescription",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Describe this in detail"}},
			},
			{
				Name:  "rsdDocuments",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "What RSD documents do you have?"}},
				Options: []i18n.Strings{
					{{"en", "Family Attestation"}},
					{{"en", "Refugee ID"}},
					{{"en", "Asylum certificate"}},
					{{"en", "Rejection decision"}},
					{{"en", "Other"}},
				},
				Multiple: true,
			},
			{
				Name:  "rsdComment",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Comment"}},
			},
			{
				Name:  "rsdIssues",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Specific RSD issues presented (narrative)"}},
			},
			//  subsection HLP
			{
				Name:  "specificHLPConcern",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "What specific HLP concern is presented?"}},
				Options: []i18n.Strings{
					{{"en", "Housing"}},
					{{"en", "Land"}},
					{{"en", "Property"}},
				},
			},
			//  sub-subsection House
			{
				Name:  "indStay",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Does the individual stay in their own house or rent?"}},
				Options: []i18n.Strings{
					{{"en", "Own house"}},
					{{"en", "Rent"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "hasRentalAgreement",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "In the case of rent, does the individual possess any agreement?"}},
			},
			{
				Name:  "agreementKind",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "What kind of agreement of proof does the individual possess?"}},
			},
			{
				Name:  "hasEvictionRisk",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Have you been or are you at risk of eviction?"}},
			},
			{
				Name:  "evictionDoc",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "If yes, what eviction document or proof do you possess?"}},
				Options: []i18n.Strings{
					{{"en", "Eviction notice"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "evictionComment",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Comment/Narrative"}},
			},
			//  sub-subsection Land
			{
				Name:  "isLegalOwner",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Are you the legal owner of the land?"}},
			},
			{
				Name:  "natureLandTenancy",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Nature of tenancy"}},
				Options: []i18n.Strings{
					{{"en", "Joint ownership"}},
					{{"en", "Co-ownership"}},
					{{"en", "Individual ownership"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "natureLandTenure",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Nature of tenure"}},
				Options: []i18n.Strings{
					{{"en", "Mailo"}},
					{{"en", "Lease"}},
					{{"en", "Freehold"}},
					{{"en", "Sustomary"}},
				},
			},
			{
				Name:  "hasLandEvictionRisk",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Have you been or are you at risk of eviction?"}},
			},
			{
				Name:  "landEvictionProof",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "If yes, what eviction document or proof do you possess?"}},
				Options: []i18n.Strings{
					{{"en", "Eviction notice"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "specificLandIssues",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Specific land issues (Narrative)"}},
			},
			// sub-subsection Property
			{
				Name:  "propertyNature",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Nature of the property in contest."}},
			},
			{
				Name:  "hasLegalOwnershipOfProperty",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Do you have legal ownership of the property?"}},
			},
			{
				Name:  "propertyOwnershipProof",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Proof of property ownership (supporting documents)"}},
			},
			{
				Name:  "propertyAcquisitionInquiry",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Inquiry on property acquisition (Narrative)"}},
			},
			// sub-section LCD
			{
				Name:  "documentationChallenges",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What documentation challenges do you have?"}},
			},
			{
				Name:  "typeOfDocument",
				Type:  form.Radio,
				Label: i18n.Strings{{"en", "What type of document do you have?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Legal"}},
						Value: "Legal",
					},
					{
						Label: i18n.Strings{{"en", "Civil"}},
						Value: "Civil",
					},
				},
			},
			{
				Name:  "lcdActionsTaken",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What action has been taken so far on this issue?"}},
			},
			// sub-section ELP
			{
				Name:  "employmentOrBusinessChallenge",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Is it and employment or business challenge?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Employment"}},
						Value: "employment",
					},
					{
						Label: i18n.Strings{{"en", "Business"}},
						Value: "business",
					},
				},
			},
			{
				Name:  "employmentChallenges",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What employment challenges do you have?"}},
			},
			{
				Name:  "typeOfEmploymentAgreement",
				Type:  form.Radio,
				Label: i18n.Strings{{"en", "What type of agreement do you have?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Oral"}},
						Value: "Legal",
					},
					{
						Label: i18n.Strings{{"en", "Written"}},
						Value: "Civil",
					},
				},
			},
			{
				Name:  "elpActionsTaken",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What actions have been taken? (Narrative)"}},
			},
			{
				Name:  "businessRelatedChallenges",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What business-related challenges do you have? (Narrative)"}},
			},
			{
				Name:  "businessRegistrationNeeds",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What business registration needs do you have? (Narrative)"}},
			},
		},
	}
	UGProtectionReferral = form.Form{
		Controls: []form.Control{
			{
				Name:  "priority",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Priority"}},
				Options: []i18n.Strings{
					{{"en", "Low (follow-up within 7 days)"}},
					{{"en", "Medium (follow-up in 3 days)"}},
					{{"en", "High (follow-up within 1 day)"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "referredVia",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Referred via"}},
				Options: []i18n.Strings{
					{{"en", "Phone (High priority only)"}},
					{{"en", "Email"}},
					{{"en", "In person"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "referralDate",
				Label: i18n.Strings{{"en", "Referral date"}},
				Type:  form.Date,
			},
			{
				Name:  "receivingAgency",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Receiving Agency"}},
			},
			{
				Name:  "partnerName",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Name of partner case worker"}},
			},
			{
				Name:  "recipientPosition",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Position of person receiving referral"}},
			},
			{
				Name:  "recipientContact",
				Type:  form.Phone,
				Label: i18n.Strings{{"en", "Contact of person receiving referral"}},
			},
			{
				Name:  "releaseConsent",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Consent to release information"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Yes"}},
						Value: "yes",
					},
					{
						Label: i18n.Strings{{"en", "No"}},
						Value: "no",
					},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "referralRestriction",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Has person expressed any restriction on referrals? If yes, specify."}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Yes"}},
						Value: "yes",
					},
				},
			},
			{
				Name:  "specification",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Specification of restriction on referrals"}},
			},
			{
				Name:  "isMinor",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Is a beneficiary a minor?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Yes"}},
						Value: "yes",
					},
					{
						Label: i18n.Strings{{"en", "No"}},
						Value: "no",
					},
				},
			},
			{
				Name:  "primaryGiver",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Name of the primary giver"}},
			},
			{
				Name:  "relationshipToChild",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Relationship to the child"}},
			},
			{
				Name:  "careGiverInformed",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Is care giver informed of referral?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Yes"}},
						Value: "yes",
					},
				},
			},
			{
				Name:  "noInformationExplanation",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "If not informed, explain"}},
			},
			{
				Name:  "referralReason",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Reason for referral"}},
			},
			{
				Name:  "referralType",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Type of referral"}},
				Options: []i18n.Strings{
					{{"en", "Health"}},
					{{"en", "Livelihood/IGAS"}},
					{{"en", "Psychosocial support"}},
					{{"en", "Safety and security"}},
					{{"en", "Education"}},
					{{"en", "Shelter"}},
				},
			},
		},
	}

	UGICLAAppointment = form.Form{
		Controls: []form.Control{
			{
				Name:  "name",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Name"}},
			},
			{
				Name:  "place",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Place"}},
			},
			{
				Name:  "date",
				Type:  form.Date,
				Label: i18n.Strings{{"en", "Date"}},
			},
			{
				Name:  "preferredContactMethod",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Preferred Contact Method"}},
				Options: []i18n.Strings{
					{{"en", "Email"}},
					{{"en", "Telephone"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "appointmentPurpose",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Appointment purpose"}},
				Options: []i18n.Strings{
					{{"en", "HLP"}},
					{{"en", "LCD"}},
					{{"en", "RSD"}},
					{{"en", "Employment/Business"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "preferredDate",
				Type:  form.Date,
				Label: i18n.Strings{{"en", "Preferred date"}},
			},
		},
	}
	UGICLAConsent = form.Form{
		Controls: []form.Control{
			{
				Name:  "consentGiven",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Has the beneficiary consented?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Consent given"}},
						Value: "yes",
					},
				},
			},
			{
				Name:  "consentProofURL",
				Type:  form.URL,
				Label: i18n.Strings{{"en", "URL to proof of beneficiary consent."}},
			},
		},
	}

	UGProtectionIncident = form.Form{
		Controls: []form.Control{
			{
				Name:  "locationOfIncident",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Location of incident"}},
			},
			{
				Name:  "timeOfIncident",
				Type:  form.Time,
				Label: i18n.Strings{{"en", "Time of incident"}},
			},
			{
				Name:  "reportedIncidentDate",
				Type:  form.Date,
				Label: i18n.Strings{{"en", "Date incident reported"}},
			},
			{
				Name:  "receivedBy",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Received by"}},
			},
			{
				Name:  "vulnerability",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Vulnerability"}},
				Options: []i18n.Strings{
					{{"en", "Child at Risk"}},
					{{"en", "Elder at Risk"}},
					{{"en", "Single parent"}},
					{{"en", "Separated Child"}},
					{{"en", "Disability"}},
					{{"en", "Woman at Risk"}},
					{{"en", "Legal and physical protection"}},
					{{"en", "Medical condition"}},
					{{"en", "Pregnant/ lactating woman"}},
				},
			},
			{
				Name:  "incidentDescription",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Description of the Incident i.e Where, when, what, who involved"}},
			},
			{
				Name:  "incidentHasBeenReportedToPolice",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Has the incident been reported to police?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Has been reported?"}},
						Value: "yes",
					},
				},
			},
			{
				Name:  "comment",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Comment"}},
			},
			{
				Name:  "incidentHasBeenReportedToOthers",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Has the Incident been reported to:"}},
				Options: []i18n.Strings{
					{{"en", "UNCHR"}},
					{{"en", "Other platforms"}},
				},
				Multiple: true,
			},
		},
	}

	UGProtectionActionReport = form.Form{
		Controls: []form.Control{
			{
				Name:  "agreedUponService",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Which service has the beneficiary together with staff agreed to take?"}},
				Options: []i18n.Strings{
					{{"en", "Cash support"}},
					{{"en", "Referral"}},
					{{"en", "Other (Specify with a narrative)"}},
					{{"en", "Relocation"}},
					{{"en", "Livelihood"}},
					{{"en", "Business support"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "narrative",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Narrative/ Describe the response action agreed."}},
			},
			{
				Name:  "followUp",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Agreed follow up with beneficiary."}},
			},
		},
	}

	UGProtectionIntake = form.Form{
		Controls: []form.Control{
			{
				Name:  "screeningDate",
				Type:  form.Date,
				Label: i18n.Strings{{"en", "Date of screening"}},
			},
			{
				Name:  "riskExposure",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Have you been exposed to any protection risk?"}},
				Options: []i18n.Strings{
					{{"en", "Violence"}},
					{{"en", "Coercion"}},
					{{"en", "Discrimination"}},
					{{"en", "Deprivation"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "riskExposureType",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "What type of protection concern experienced?"}},
				Options: []i18n.Strings{
					{{"en", "Physical violence"}},
					{{"en", "Neglect"}},
					{{"en", "Family separation"}},
					{{"en", "Arrest"}},
					{{"en", "Denial of resources"}},
					{{"en", "Psychosocial violence"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "details",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Provide details (Narrative)"}},
			},
			{
				Name:  "responsePriority",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Response Priority"}},
				Options: []i18n.Strings{
					{{"en", "High (follow up requested in 24 hours)"}},
					{{"en", "Medium (Follow up in 3 days)"}},
					{{"en", "Low (Follow up in 7 days)"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
		},
	}

	UGProtectionFollowUp = form.Form{
		Controls: []form.Control{
			{
				Name:  "followUpAfter",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Follow up after"}},
				Options: []i18n.Strings{
					{{"en", "1 week"}},
					{{"en", "2 weeks"}},
					{{"en", "1 month"}},
					{{"en", "3 months"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "agreedFollowUp",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Agreed follow up with the beneficiary"}},
			},
		},
	}

	UGICLAActionPlan = form.Form{
		Controls: []form.Control{
			{
				Name:     "agreedService",
				Type:     form.Dropdown,
				Label:    i18n.Strings{{"en", "Which service has the beneficiary together with staff agreed to take?"}},
				Multiple: true,
				Options: []i18n.Strings{
					{{"en", "Legal counselling"}},
					{{"en", "Referral"}},
					{{"en", "Other (Specify with a narrative),"}},
					{{"en", "Relocation"}},
					{{"en", "Livelihood"}},
					{{"en", "Business support"}},
				},
			},
			{
				Name:     "agreedAction",
				Type:     form.Dropdown,
				Label:    i18n.Strings{{"en", "Type of actions for case worker agreed upon with beneficiary"}},
				Multiple: true,
				Options: []i18n.Strings{
					{{"en", "(1) Discussion with supervisor/team leader"}},
					{{"en", "(2) Conducting legal analysis, including the study of judicial practice"}},
					{{"en", "(3) Preparing letters, inquiries to various authorities"}},
					{{"en", "(4) Drafting of other legal documents (such leases or contracts)"}},
					{{"en", "(5) Lodging of a court application"}},
					{{"en", "(6) Attending of court session/hearing"}},
					{{"en", "(7) Review of the decision/appeal"}},
					{{"en", "(8) Negotiation"}},
					{{"en", "(9) Follow up with relevant administrative authority or other entities"}},
					{{"en", "(10) Accompaniment"}},
					{{"en", "(11) Other"}},
				},
			},
			{
				Name:  "actionComment",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Action comment"}},
			},
		},
	}

	specificNeedsOptions = []i18n.Strings{
		{{"en", "Disability"}},
		{{"en", "Pregant woman"}},
		{{"en", "Elderly taking care of minors alone"}},
		{{"en", "Elderly living alone"}},
		{{"en", "Single parent"}},
		{{"en", "Chronic illness"}},
		{{"en", "Legal protection needs"}},
		{{"en", "Child"}},
		{{"en", "Other"}},
	}

	disabilityNeedsOptions = []i18n.Strings{
		{{"en", "No"}},
		{{"en", "Moderate physical impairment"}},
		{{"en", "Severe physical impairment"}},
		{{"en", "Moderate sensory impairment"}},
		{{"en", "Severe sensory impairment"}},
		{{"en", "Moderate mental disability"}},
		{{"en", "Severe mental disability"}},
	}

	obstacleOptions = []i18n.Strings{
		{{"en", "Insufficient funds"}},
		{{"en", "Distance issues"}},
		{{"en", "Insecurity"}},
		{{"en", "Social discrimination"}},
		{{"en", "Insufficient quantity of goods"}},
		{{"en", "Inadequate quality of goods/services"}},
		{{"en", "Insufficient capabilities and competences"}},
		{{"en", "Others"}},
	}

	meetNeedsAbility = []i18n.Strings{
		{{"en", "Can meet all needs wont worry"}},
		{{"en", "Can meet needs"}},
		{{"en", "Can barely meet needs"}},
		{{"en", "Unable to meet needs"}},
		{{"en", "Totally unable to meet needs"}},
	}

	UGProtectionSocialStatusAssessment = form.Form{
		Controls: []form.Control{
			{
				Name:    "specificNeeds",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Does the Client have specific needs?"}},
				Options: specificNeedsOptions,
			},
			{
				Name:  "comment",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Comment"}},
			},
			{
				Name:  "otherHHMemberHasSpecificNeeds",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Does any other member of the HH have specific needs?"}},
			},
			{
				Name:    "otherHHMemberSpecificNeeds",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "What specific need does the HH member have?"}},
				Options: specificNeedsOptions,
			},
			{
				Name:  "homeSituation",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Home situation"}},
				Options: []i18n.Strings{
					{{"en", "Lives alone"}},
					{{"en", "Lives with family"}},
					{{"en", "Hosted by relatives"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:    "disability",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Does the client have a disability?"}},
				Options: disabilityNeedsOptions,
			},
			{
				Name:    "otherHHMemberDisability",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Does any other member of the HH live with disability?"}},
				Options: disabilityNeedsOptions,
			},
			{
				Name:  "disabledHHMember",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Which other member of the HH lives with a disability?"}},
			},
			{
				Name:  "incomeNeeds",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "How do you meet your income HH needs for all its members?"}},
				Options: []i18n.Strings{
					{{"en", "Remittances"}},
					{{"en", "Savings and credit groups"}},
					{{"en", "Small business(registered)"}},
					{{"en", "Small business (unregistered)"}},
					{{"en", "Unskilled labor"}},
					{{"en", "Skilled labor"}},
					{{"en", "Agriculture/pastoralism"}},
					{{"en", "Donation/Humanitarian assistance"}},
					{{"en", "Begging"}},
					{{"en", "Bank loan"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "workingMembersOfHH",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "How many people are able to work in your HH?"}},
			},
			{
				Name:  "workAmount",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "How often do they work?"}},
			},
			{
				Name:  "workType",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What do they do?"}},
			},
			{
				Name:  "humanitarianAssistance",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Do you receive humanitarian assistance?"}},
			},
			{
				Name:  "humanitarianAssistanceComment",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Comment/ recent changes regarding humanitarian assistance"}},
			},
			{
				Name:  "homeSituationComment",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Comment on living situation (Narrative)"}},
			},
			{
				Name:  "amountGirls0to5",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 0-5 old girls"}},
			},
			{
				Name:  "amountBoys0to5",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 0-5 old boys"}},
			},
			{
				Name:  "amountGirls6to12",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 6-12 old girls"}},
			},
			{
				Name:  "amountBoys6to12",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 6-12 old boys"}},
			},
			{
				Name:  "amountGirls13to17",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 13-17 old girls"}},
			},
			{
				Name:  "amountBoys13to17",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 13-17 old boys"}},
			},
			{
				Name:  "amountFemales18to59",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 18-59 old females"}},
			},
			{
				Name:  "amountMales18to59",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 18-59 old males"}},
			},
			{
				Name:  "amountFemalesOver59",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 59+ old females"}},
			},
			{
				Name:  "amountMalesOver59",
				Type:  form.Number,
				Label: i18n.Strings{{"en", "Number of 59+ old males"}},
			},
			{
				Name:    "foodNeedsMet",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "HH’s ability to meet the food needs of all its members."}},
				Options: meetNeedsAbility,
			},
			{
				Name:    "foodNeedObstacles",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "What are the main obstacles you face in meeting food needs?"}},
				Options: obstacleOptions,
			},
			{
				Name:    "accommodationNeedObstacles",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Main Obstacles you face in meeting accommodation needs"}},
				Options: obstacleOptions,
			},
			{
				Name:    "washNeedsMet",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Can the HH meet WASH needs?"}},
				Options: meetNeedsAbility,
			},
			{
				Name:    "washNeedsObstacles",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", " Main obstacles in meeting WASH needs?"}},
				Options: obstacleOptions,
			},
			{
				Name:  "summaryNarrative",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Summary Narrative"}},
			},
		},
	}

	UGICLAReferral = form.Form{
		Sections: []form.Section{
			{
				Title: []i18n.LocaleString{{"en", "Beneficiary's Information"}},
				ControlNames: []string{
					"typeOfICLAReferral",
					"typeOfServicesRequested",
					"reasonForReferral",
				},
			},
			{
				Title: []i18n.LocaleString{{"en", "Consent"}},
				ControlNames: []string{
					"beneficiaryConsentedToInformationRelease",
					"beneficiaryConsentedToInformationReleaseHintYes",
					"beneficiaryConsentedToInformationReleaseHintNo",
				},
			},
			{
				Title: i18n.Strings{{"en", "Means of Referral"}},
				ControlNames: []string{
					"beneficiaryHasRestrictionsWithReferral",
					"meansOfReferral",
					"meansOfFeedback",
				},
			},
		},
		Controls: []form.Control{
			{
				Name:  "typeOfICLAReferral",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Type of Referral"}},
				Options: []i18n.Strings{
					{{"en", "Internal - Shelter/NFI"}},
					{{"en", "Internal - Livelihood/Food Security"}},
					{{"en", "Internal - Education"}},
					{{"en", "Internal - WASH"}},
					{{"en", "Internal - Camp Management/UDOC"}},
					{{"en", "External - Organisation"}},
					{{"en", "External - Contact person"}},
					{{"en", "External - Phone number"}},
					{{"en", "External - E-Mail"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:     "typeOfServicesRequested",
				Type:     form.Dropdown,
				Label:    i18n.Strings{{"en", "Type of Services/Assistance Requested"}},
				Multiple: true,
				Options: []i18n.Strings{
					{{"en", "Health care (including medication)"}},
					{{"en", "Legal assistance"}},
					{{"en", "Education"}},
					{{"en", "Mental Health"}},
					{{"en", "Transportation"}},
					{{"en", "Food"}},
					{{"en", "Non-food items (including hygiene items)"}},
					{{"en", "Disability"}},
					{{"en", "MPC"}},
					{{"en", "Shelter/Housing"}},
					{{"en", "Shelter construction/repair"}},
					{{"en", "Youth Livelihoods (e.g. vocational training)"}},
					{{"en", "Small/Medium Business Grants"}},
					{{"en", "Other livelihood activities"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "reasonForReferral",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Reason for Referral"}},
			},
			{
				Name:  "beneficiaryConsentedToInformationRelease",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Has the beneficiary consented to the release of his/her information for the referral? "}},
			},
			{
				Name:  "beneficiaryConsentedToInformationReleaseHintYes",
				Type:  form.Hint,
				Label: i18n.Strings{{"en", "If 'yes', upload signed consent form and proceed"}},
			},
			{
				Name:  "beneficiaryConsentedToInformationReleaseHintNo",
				Type:  form.Hint,
				Label: i18n.Strings{{"en", "If 'no', explain the reason why and do not refer the case"}},
			},
			{
				Name:  "beneficiaryHasRestrictionsWithReferral",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Does the beneficiary have any restrictions to be referred?"}},
			},
			{
				Name:  "meansOfReferral",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Means of referral "}},
				Options: []i18n.Strings{
					{{"en", "(1) Phone"}},
					{{"en", "(2) E-mail"}},
					{{"en", "(3) Personal meeting"}},
					{{"en", "(4) Other"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:  "meansOfFeedback",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Means and terms of receiving feedback from the client"}},
				Options: []i18n.Strings{
					{{"en", "(1) Phone"}},
					{{"en", "(2) E-mail"}},
					{{"en", "(3) Personal meeting"}},
					{{"en", "(4) Other"}},
				},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
		},
	}

	// Case Types for Uganda
	// - Kampala Response Team
	UGSituationalAnalysisCaseType = caseType("0ae90b08-6944-48dc-8f30-5cb325292a8c", "Situational Analysis (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGSituationAnalysis, true)
	UGIndividualResponseCaseType  = caseType("2f909038-0ce4-437b-af17-72fc5d668b49", "Response (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGIndividualResponse, true)
	//UGReferralCaseType                 = caseType("ecdaf47f-6fa9-48c8-9d10-6324bf932ed7", "Referral (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGReferral, false)
	//UGExternalReferralFollowupCaseType = caseType("2a1b670c-6336-4364-b89d-0e65fc771659", "External Referral Followup (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGExternalReferralFollowup, false)
	UGProtectionIntakeCaseType                 = caseType("da20a49d-3cc9-413c-89b8-ff40e3afe95c", "Intake (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionIntake, true)
	UGProtectionFollowUpCaseType               = caseType("dcebe6c8-47cd-4e0f-8562-5680573aed88", "Follow up (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionFollowUp, false)
	UGProtectionSocialStatusAssessmentCaseType = caseType("e3b30f91-7181-41a3-8187-f176084a0ab2", "Social Status Assessment (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionSocialStatusAssessment, false)
	UGProtectionReferralCaseType               = caseType("dc18bf9d-e812-43a8-b843-604c23306cd6", "UG Protection Referral (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionReferral, false)
	UGProtectionIncidentCaseType               = caseType("f6117a29-db5a-49d7-b564-bf42740ae824", "Incident (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionIncident, false)
	UGProtectionActionReportCaseType           = caseType("f4989460-8e76-4d82-aad5-ed2ad3d3d627", "Action Report (UG Protection/Response)", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionActionReport, false)

	// - Kampala ICLA Team
	UGICLAFollowUpCaseType       = caseType("415be6d4-cf1b-484a-9bad-83acd8474498", "ICLA Follow up (UG ICLA)", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAFollowUp, false)
	UGICLAIntakeCaseType         = caseType("61fb6d03-2374-4bea-9374-48fc10500f81", "ICLA Intake (UG ICLA)", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAIntake, true)
	UGICLACaseAssessmentCaseType = caseType("bbf820de-8d10-49eb-b8c9-728993ab0b73", "ICLA Case Assessment (UG ICLA)", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLACaseAssessment, false)
	UGICLAAppointmentCaseType    = caseType("27064ded-fbfe-4197-830c-164a797d5306", "ICLA Appointment (UG ICLA)", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAAppointment, false)
	UGICLAConsentCaseType        = caseType("3ad2d524-4dd0-4834-9fc2-47808cf66941", "ICLA Consent (UG ICLA)", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAConsent, false)
	UGICLAReferralCaseType       = caseType("9896c0f1-8d66-4657-92f2-e67a7afcf9ab", "ICLA Referral (UG ICLA)", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAReferral, false)
	UGICLAActionPlanCaseType     = caseType("2b4f46a7-aebd-4754-89fd-dc7897a79ddb", "ICLA Action Plan (UG ICLA)", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAActionPlan, false)

	// Registration Controller Flow for Uganda Intake Process
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

	// Individuals - UG Beneficiaries
	JohnDoe     = ugandaIndividual(individual("c529d679-3bb6-4a20-8f06-c096f4d9adc1", "John Sinclair Doe", "John Doe", "1983-04-23", "john.doe", "Refugee", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "No", "Yes", "Moderate", "No", "", "No", "", "Kenya", "Kiswahili, English", "English", "123 Main Street, Kampala", "0123456789", "", "Email", "No"), "1983-04-23", "0", "0", "0", "0", "0", "0")
	MaryPoppins = ugandaIndividual(individual("bbf539fd-ebaa-4438-ae4f-8aca8b327f42", "Mary Poppins", "Mary Poppins", "1983-04-23", "mary.poppins", "Internally Displaced Person", "Female", "Yes", "https://link-to-consent.proof", "No", "No", "No", "No", "", "No", "", "No", "", "Uganda", "Rukiga, English", "Rukiga", "901 First Avenue, Kampala", "0123456789", "", "Telegram", "Yes"), "1983-04-23", "0", "0", "0", "0", "0", "0")
	BoDiddley   = ugandaIndividual(individual("26335292-c839-48b6-8ad5-81271ee51e7b", "Ellas McDaniel", "Bo Diddley", "1983-04-23", "bo.diddley", "Host Community", "Male", "Yes", "https://link-to-consent.proof", "No", "No", "Yes", "No", "", "No", "", "No", "", "Somalia", "Somali, Arabic, English", "English", "101 Main Street, Kampala", "0123456789", "", "Whatsapp", "No"), "1983-04-23", "0", "0", "0", "0", "0", "0")

	// Individuals - UG Staff
	Stephen   = individual("066a0268-fdc6-495a-9e4b-d60cfae2d81a", "Stephen Kabagambe", "Stephen Kabagambe", "1983-04-23", "stephen.kabagambe", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Colette   = individual("93f9461f-31da-402e-8988-6e0100ecaa24", "Colette le Jeune", "Colette le Jeune", "1983-04-23", "colette.le.jeune", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Courtney  = individual("14c014d9-f433-4508-b33d-dc45bf86690b", "Courtney Lare", "Courtney Lare", "1983-04-23", "courtney.lare", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Ulrika    = individual("8f3abe4a-f6c2-45df-a095-482eb4f9a3e9", "Ulrika Blom", "Ulrika Blom", "1983-04-23", "ulrika.blom", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Nathan    = individual("8fd700e7-3621-4b9c-8f25-db0a90994002", "Nathan Chelimo", "Nathan Chelimo", "1983-04-23", "nathan.chelimo", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Lilian    = individual("4db00c54-c15a-4594-932c-08127a35e6c8", "Lilian Kabooga", "Lilian Kabooga", "1983-04-23", "lilian.kabooga", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Viviane   = individual("4517afc4-d105-4585-9f6f-9b2643ba03b9", "Viviane Mushimiyimana", "Viviane Mushimiyimana", "1983-04-23", "viviane.n", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Elizabeth = individual("ee1002bf-a8cf-4c29-95e2-57dbe1ef3054", "Elizabeth Salaama", "Elizabeth Salaama", "1983-04-23", "elizabeth.salaama", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Basil     = individual("72bf35ab-1198-46c5-bc8d-b86a9b23a6fa", "Basil Droti", "Basil Droti", "1983-04-23", "basil.droti", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Abaasi    = individual("0d816f2f-51d4-4b78-b4fb-7d65f5712111", "Abaasi Abdilahin", "Abaasi Abdilahin", "1983-04-23", "abaasi.abdilahin", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Judith    = individual("69ff82af-ecdc-49ca-9f2e-a90d3eca706a", "Judith Andoua", "Judith Andoua", "1983-04-23", "judith.andoua", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Muriel    = individual("1923aaea-8b82-4617-9d41-c345a5235bb4", "Muriel Iyanu", "Muriel Iyanu", "1983-04-23", "muriel.iyanu", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Christine = individual("c42f8aa6-e873-48d7-8f51-8bb312ea6ae1", "Christine Onen", "Christine Onen", "1983-04-23", "christine.onen", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Jamusi    = individual("5f1dbd2d-62e0-4f25-97f7-f982da6553d6", "Jamusi Kisyenene", "Jamusi Kisyenene", "1983-04-23", "jamusi.kisyenene", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Juliet    = individual("07eae976-fd58-415a-bb91-8b19de7ba5fc", "Juliet Aryamo", "Juliet Aryamo", "1983-04-23", "juliet.aryamo", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Rebecca   = individual("3883a43f-95af-4a3a-9490-7f726a76f169", "Rebecca Naluzze", "Rebecca Naluzze", "1983-04-23", "rebecca.naluzze", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Angella   = individual("9096ee30-2fd8-4c7d-a038-518a6e2e6b44", "Angella Namuyomba", "Angella Namuyomba", "1983-04-23", "angella.namuyomba", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Hassan    = individual("2bab647f-1e1e-46a7-a721-125c9214d345", "Hassan Mpanga", "Hassan Mpanga", "1983-04-23", "hassan.mpanga", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Fauzia    = individual("f159b04d-f7d2-4d4e-9ff4-e267a1caf566", "Fauzia Nkunyingi", "Fauzia Nkunyingi", "1983-04-23", "fauzia.nkunyingi", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Denis     = individual("679ec5ed-1cb4-4a4c-b872-0e7dd97ab7c5", "Denis Onena", "Denis Onena", "1983-04-23", "denis.onena", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Nicola    = individual("fef9c068-5603-4469-9d9a-3b48706d2424", "Nicola Cozza", "Nicola Cozza", "1983-04-23", "nicola.cozza", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Grace     = individual("bf672aa8-a22f-47e2-ae10-1c72f12bdd90", "Grace Chebet", "Grace Chebet", "1983-04-23", "grace.chebet", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Robert    = individual("8b645f9f-d727-4bad-830c-a04a2f8dd0f1", "Robert Dikua", "Robert Dikua", "1983-04-23", "robert.dikua", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Emmaneul  = individual("0c57d265-2c08-4b09-813d-afd62184f878", "Emmaneul Epaire", "Emmaneul Epaire", "1983-04-23", "emmaneul.epaire", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Shamir    = individual("11b3394b-d621-4b90-8f3c-5ff70e57c48a", "Shamir Nabada", "Shamir Nabada", "1983-04-23", "shamir.nabada", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")

	// Individuals - CO staff
	Claudia = individual("0888928f-aa48-4b5f-a23e-8f885d734f71", "Claudia Garcia", "Claudia Garcia", "1983-04-23", "claudia.garcia", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")

	// Beneficiaries
	_ = beneficiary(JohnDoe)
	_ = beneficiary(MaryPoppins)
	_ = beneficiary(BoDiddley)

	// Staff
	_ = staff(Stephen)
	_ = staff(Colette)
	_ = staff(Courtney)
	_ = staff(Claudia)
	_ = staff(Ulrika)
	_ = staff(Nathan)
	_ = staff(Lilian)
	_ = staff(Viviane)
	_ = staff(Elizabeth)
	_ = staff(Basil)
	_ = staff(Abaasi)
	_ = staff(Judith)
	_ = staff(Muriel)
	_ = staff(Christine)
	_ = staff(Jamusi)
	_ = staff(Juliet)
	_ = staff(Rebecca)
	_ = staff(Angella)
	_ = staff(Hassan)
	_ = staff(Fauzia)
	_ = staff(Denis)
	_ = staff(Nicola)
	_ = staff(Grace)
	_ = staff(Robert)
	_ = staff(Emmaneul)
	_ = staff(Shamir)

	// Memberships
	StevenMembership    = membership("862690ee-87f0-4f95-aa1e-8f8a2f2fd54a", Stephen, KampalaCOTeam)
	ColetteMembership   = membership("9d4abef9-0be0-4750-81ab-0524a412c049", Colette, KampalaProtectionTeam)
	CourtneyMembership  = membership("83c5e73a-5947-4d7e-996c-14a2a7b1c850", Courtney, KampalaProtectionTeam)
	ClaudiaMembership   = membership("344016e3-1d89-4f28-976b-1bf891d69aff", Claudia, ColombiaTeam)
	UlrikaMembership    = membership("06e69c81-3c43-4c6d-b950-723af73ad5aa", Ulrika, KampalaCOTeam)
	NathanMembership    = membership("32d48a07-9e4c-4e3c-b0f2-68e821b013a0", Nathan, KampalaCOTeam)
	LilianMembership    = membership("fd781438-7ecf-412a-b5a7-2b0db67d06f2", Lilian, KampalaProtectionTeam)
	VivianeMembership   = membership("cce94aac-4dc6-499f-9333-bdf2db9a1e7e", Viviane, KampalaProtectionTeam)
	ElizabethMembership = membership("7e86fbcc-f1c3-46dc-b363-d241ce74ac67", Elizabeth, KampalaProtectionTeam)
	BasilMembership     = membership("d0a68559-ebbb-46fd-97ba-b88ce51995a0", Basil, KampalaICLATeam)
	AbaasiMembership    = membership("c92f4e4a-5c44-4832-9a8c-0d676f168685", Abaasi, KampalaICLATeam)
	JudithMembership    = membership("2f9dfba9-e2cb-4634-8616-dcc0b7de68ae", Judith, KampalaICLATeam)
	MurielMembership    = membership("302a21f8-692e-422d-8f33-67b1c124f481", Muriel, KampalaICLATeam)
	ChristineMembership = membership("bd69517c-e4e7-4e4d-a19c-a12d605fb830", Christine, KampalaCOTeam)
	JamusiMembership    = membership("8de2e1fc-6d04-430b-badd-0726fcc3e006", Jamusi, KampalaCOTeam)
	JulietMembership    = membership("ae0b2a4d-4727-4efc-a1c2-8d6ecdabd062", Juliet, KampalaCOTeam)
	RebeccaMembership   = membership("f8543015-4634-4115-a9c0-1c8df56ec0ec", Rebecca, KampalaCOTeam)
	AngellaMembership   = membership("eb97df7a-978b-413d-998c-4d14269f59ea", Angella, KampalaCOTeam)
	HassanMembership    = membership("b78f7bed-7b28-4025-90e1-b9f2aaaf332c", Hassan, KampalaCOTeam)
	FauziaMembership    = membership("7f5e9a79-0533-439c-adbd-5b4e411f538f", Fauzia, KampalaCOTeam)
	DenisMembership     = membership("dc42a6d1-635a-402a-8376-638dd263c549", Denis, KampalaCOTeam)
	NicolaMembership    = membership("d30188f6-12a5-4485-85f0-a42be0247e3c", Nicola, KampalaCOTeam)
	GraceMembership     = membership("ced861ed-120b-4af9-98ce-aa1000fef05e", Grace, KampalaCOTeam)
	RobertMembership    = membership("9c438611-b11e-4a76-9e44-ac14ff3e14c6", Robert, KampalaCOTeam)
	EmmaneulMembership  = membership("ce8f6666-feea-40ea-9c71-d92df73f5720", Emmaneul, KampalaProtectionTeam)
	ShamirMembership    = membership("c5776d66-ca35-483d-bcf3-721fcde9eeff", Shamir, KampalaProtectionTeam)

	// Countries
	ugandaCountry   = country(iam.UgandaCountry.ID, iam.UgandaCountry.Name)
	colombiaCountry = country(iam.ColombiaCountry.ID, iam.ColombiaCountry.Name)

	// Nationalities
	KampalaCOTeamNationality         = nationality("0987460d-c906-43cd-b7fd-5e7afca0d93e", KampalaCOTeam, ugandaCountry)
	KampalaProtectionTeamNationality = nationality("b58e4d26-fe8e-4442-8449-7ec4ca3d9066", KampalaProtectionTeam, ugandaCountry)
	KampalaICLATeamNationality       = nationality("23e3eb5e-592e-42e2-8bbf-ee097d93034c", KampalaICLATeam, ugandaCountry)
	ColombiaTeamNationality          = nationality("7ba6d2ee-1af9-447c-8000-7719467b3414", ColombiaTeam, colombiaCountry)

	// Cases

	BoDiddleySituationAnalysisData = map[string][]string{
		"safeDignifiedLife":  {"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life."},
		"challengesBarriers": {"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life."},
		"solutions":          {"A qualified interpreter, who knows the legal context could help us to agree on contractual matters."},
		"workTogether":       {"NRC could provide a translator and a legal representative to ease contract negotiations"},
	}

	BoDiddleyResponseData = map[string][]string{
		"servicesStartingPoint": {"ICLA"},
		"commentStartingPoint":  {"The individual has requested ICLA as a starting point, we should create a referral"},
		"otherServices":         {"Protection"},
		"commentOtherServices":  {"The individual has requested additional Protection services, we should create a referral"},
		"perceivedPriority":     {"High"},
	}
	MaryPoppinsSituationAnalysisData = map[string][]string{
		"safeDignifiedLife":  {"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life."},
		"challengesBarriers": {"Some of the barriers I face are communication gaps between myself and refugee tenants. We are attempting to deal with these challenges by using google translate."},
		"solutions":          {"A qualified interpreter, who knows the legal context could help us to agree on contractual matters."},
		"workTogether":       {"NRC could provide a translator and a legal representative to ease contract negotiations"},
	}
	MaryPoppinsResponseData = map[string][]string{
		"servicesStartingPoint": {"S&S"},
		"commentStartingPoint":  {"The individual has requested S&S as a starting point, we should create a referral"},
		"otherServices":         {"Protection"},
		"commentOtherServices":  {"The individual has requested additional Protection services, we should create a referral"},
		"perceivedPriority":     {"High"},
	}
	JohnDoeSituationAnalysisData = map[string][]string{
		"safeDignifiedLife":  {"Yes, I live a safe and dignified life and I am reasonably happy with my achievements and quality of life."},
		"challengesBarriers": {"Some of the barriers I face are communication gaps between myself and refugee tenants. We are attempting to deal with these challenges by using google translate."},
		"solutions":          {"A qualified interpreter, who knows the legal context could help us to agree on contractual matters."},
		"workTogether":       {"NRC could provide a translator and a legal representative to ease contract negotiations"},
	}
	JohnDoeResponseData = map[string][]string{
		"servicesStartingPoint": {"LFS"},
		"commentStartingPoint":  {"The individual has requested LFS as a starting point, we should create a referral"},
		"otherServices":         {"WASH"},
		"commentOtherServices":  {"The individual has requested additional WASH services, we should create a referral"},
		"perceivedPriority":     {"High"},
	}

	BoDiddleySituationAnalysis  = kase("dba43642-8093-4685-a197-f8848d4cbaaa", Colette.ID, BoDiddley.ID, KampalaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, BoDiddleySituationAnalysisData)
	BoDiddleyIndividualResponse = kase("3ea8c121-bdf0-46a0-86a8-698dc4abc872", Colette.ID, BoDiddley.ID, KampalaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, BoDiddleyResponseData)

	MaryPoppinsSituationAnalysis  = kase("4f7708ed-240a-423f-9bd1-839542e65833", Colette.ID, MaryPoppins.ID, KampalaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, MaryPoppinsSituationAnalysisData)
	MaryPoppinsIndividualResponse = kase("45b4a637-c610-4ab9-afe6-4e958c36a96f", Colette.ID, MaryPoppins.ID, KampalaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, MaryPoppinsResponseData)

	JohnDoesSituationAnalysis = kase("43140381-8166-4fb3-9ac5-339082920ade", Colette.ID, JohnDoe.ID, KampalaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, JohnDoeSituationAnalysisData)
	JohnDoeIndividualResponse = kase("65e02e79-1676-4745-9890-582e3d67d13f", Colette.ID, JohnDoe.ID, KampalaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, JohnDoeResponseData)

	// Identification Document Types
	DriversLicense = identificationDocumentType("75c41c5f-bf7e-4b45-a242-5e0f875e3044", "Drivers License")
	NationalID     = identificationDocumentType("8910a1ea-4bfe-4321-aa5b-15922b09ad4d", "National ID")
	UNHCRID        = identificationDocumentType("6833cb6d-593f-4f3f-926d-498be74352d1", "UNHCR ID")
	Passport       = identificationDocumentType("567d04e5-abf4-4899-848f-0395264309f0", "Passport")

	BoDiddleyPassport  = identificationDocument("20d194d6-a1ac-483e-8c24-38b5efbaca6f", BoDiddley.ID, "A0JBODIDDLEY129", Passport.ID)
	MaryPoppinsUNHRCID = identificationDocument("0244b59e-5d5c-4e13-af96-da1ccf4e9499", MaryPoppins.ID, "LLP987MARYPOPPINS99", UNHCRID.ID)
	JohnDoeNationalID  = identificationDocument("4c9477c9-c149-4db7-928c-f5e5f915e018", JohnDoe.ID, "B811HJOHNDOE01", NationalID.ID)
)
