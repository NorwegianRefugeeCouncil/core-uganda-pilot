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
	UgandaProtectionTeam = team("ac9b8d7d-d04d-4850-9a7f-3f93324c0d1e", "Uganda Protection Team")
	UgandaICLATeam       = team("a43f84d5-3f8a-48c4-a896-5fb0fcd3e42b", "Uganda ICLA Team")
	UgandaCoreAdminTeam  = team("814fc372-08a6-4e6b-809b-30ebb51cb268", "Uganda Core Admin Team")
	ColombiaTeam         = team("a6bc6436-fcea-4738-bde8-593e6480e1ad", "Colombia Team")

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
				Name:  "workType",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What do they do?"}},
			},
			{
				Name:  "workType",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "What do they do?"}},
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
				Name:  "foodNeedsMet",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "HH’s ability to meet the food needs of all its members."}},
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
				Name:  "washNeedsMet",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", "Can the HH meet WASH needs?"}},
				Options: meetNeedsAbility,
			},
			{
				Name:  "incomeNeeds",
				Type:  form.Dropdown,
				Label: i18n.Strings{{"en", " Main obstacles in meeting WASH needs?"}},
				Options: obstacleOptions,
			},
			{
				Name:  "summaryNarrative",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Summary Narrative"}},
			},
		},
	}

	// Case Types for Uganda
	// - Kampala Response Team
	UGSituationalAnalysisCaseType              = caseType("0ae90b08-6944-48dc-8f30-5cb325292a8c", "Situational Analysis (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGSituationAnalysis, true)
	UGIndividualResponseCaseType               = caseType("2f909038-0ce4-437b-af17-72fc5d668b49", "Response (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGIndividualResponse, true)
	UGReferralCaseType                         = caseType("ecdaf47f-6fa9-48c8-9d10-6324bf932ed7", "Referral (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGReferral, false)
	UGExternalReferralFollowupCaseType         = caseType("2a1b670c-6336-4364-b89d-0e65fc771659", "External Referral Followup (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGExternalReferralFollowup, false)
	UGProtectionIntakeCaseType                 = caseType("da20a49d-3cc9-413c-89b8-ff40e3afe95c", "Intake (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionIntake, true)
	UGProtectionFollowUpCaseType               = caseType("dcebe6c8-47cd-4e0f-8562-5680573aed88", "Follow up (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionFollowUp, false)
	UGProtectionSocialStatusAssessmentCaseType = caseType("e3b30f91-7181-41a3-8187-f176084a0ab2", "Social Status Assessment (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionSocialStatusAssessment, false)
	UGProtectionReferralCaseType       = caseType("dc18bf9d-e812-43a8-b843-604c23306cd6", "UG Protection Referral (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionReferral, false)
	UGProtectionIncidentCaseType       = caseType("f6117a29-db5a-49d7-b564-bf42740ae824", "Incident (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionIncident, false)
	UGProtectionActionReportCaseType   = caseType("f4989460-8e76-4d82-aad5-ed2ad3d3d627", "Action Report (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionActionReport, false)

	// - Kampala ICLA Team
	UGICLAIntakeCaseType         = caseType("61fb6d03-2374-4bea-9374-48fc10500f81", "ICLA Intake (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLAIntake, true)
	UGICLACaseAssessmentCaseType = caseType("bbf820de-8d10-49eb-b8c9-728993ab0b73", "ICLA Case Assessment (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLACaseAssessment, false)
	UGICLAAppointmentCaseType    = caseType("27064ded-fbfe-4197-830c-164a797d5306", "ICLA Appointment (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLAAppointment, false)
	UGICLAConsentCaseType        = caseType("3ad2d524-4dd0-4834-9fc2-47808cf66941", "ICLA Consent (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLAConsent, false)

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
	Stephen  = individual("066a0268-fdc6-495a-9e4b-d60cfae2d81a", "Stephen Kabagambe", "Stephen Kabagambe", "1983-04-23", "stephen.kabagambe", "", "Male", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Colette  = individual("93f9461f-31da-402e-8988-6e0100ecaa24", "Colette le Jeune", "Colette le Jeune", "1983-04-23", "colette.le.jeune", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	Courtney = individual("14c014d9-f433-4508-b33d-dc45bf86690b", "Courtney Lare", "Courtney Lare", "1983-04-23", "courtney.lare", "", "Female", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")

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

	// Memberships
	StevenMembership   = membership("862690ee-87f0-4f95-aa1e-8f8a2f2fd54a", Stephen, UgandaCoreAdminTeam)
	ColetteMembership  = membership("9d4abef9-0be0-4750-81ab-0524a412c049", Colette, UgandaProtectionTeam)
	CourtneyMembership = membership("83c5e73a-5947-4d7e-996c-14a2a7b1c850", Courtney, UgandaProtectionTeam)
	ClaudiaMembership  = membership("344016e3-1d89-4f28-976b-1bf891d69aff", Claudia, ColombiaTeam)

	// Countries
	ugandaCountry   = country(iam.UgandaCountry.ID, iam.UgandaCountry.Name)
	colombiaCountry = country(iam.ColombiaCountry.ID, iam.ColombiaCountry.Name)

	// Nationalities

	UgandaCoreAdminTeamNationality  = nationality("0987460d-c906-43cd-b7fd-5e7afca0d93e", UgandaCoreAdminTeam, ugandaCountry)
	UgandaProtectionTeamNationality = nationality("b58e4d26-fe8e-4442-8449-7ec4ca3d9066", UgandaProtectionTeam, ugandaCountry)
	UgandaICLATeamNationality       = nationality("23e3eb5e-592e-42e2-8bbf-ee097d93034c", UgandaICLATeam, ugandaCountry)
	ColombiaTeamNationality         = nationality("7ba6d2ee-1af9-447c-8000-7719467b3414", ColombiaTeam, colombiaCountry)

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

	BoDiddleySituationAnalysis  = kase("dba43642-8093-4685-a197-f8848d4cbaaa", Colette.ID, BoDiddley.ID, UgandaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, BoDiddleySituationAnalysisData)
	BoDiddleyIndividualResponse = kase("3ea8c121-bdf0-46a0-86a8-698dc4abc872", Colette.ID, BoDiddley.ID, UgandaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, BoDiddleyResponseData)

	MaryPoppinsSituationAnalysis  = kase("4f7708ed-240a-423f-9bd1-839542e65833", Colette.ID, MaryPoppins.ID, UgandaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, MaryPoppinsSituationAnalysisData)
	MaryPoppinsIndividualResponse = kase("45b4a637-c610-4ab9-afe6-4e958c36a96f", Colette.ID, MaryPoppins.ID, UgandaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, MaryPoppinsResponseData)

	JohnDoesSituationAnalysis = kase("43140381-8166-4fb3-9ac5-339082920ade", Colette.ID, JohnDoe.ID, UgandaProtectionTeam.ID, UGSituationalAnalysisCaseType, true, true, JohnDoeSituationAnalysisData)
	JohnDoeIndividualResponse = kase("65e02e79-1676-4745-9890-582e3d67d13f", Colette.ID, JohnDoe.ID, UgandaProtectionTeam.ID, UGIndividualResponseCaseType, true, true, JohnDoeResponseData)

	// Identification Document Types
	DriversLicense = identificationDocumentType("75c41c5f-bf7e-4b45-a242-5e0f875e3044", "Drivers License")
	NationalID     = identificationDocumentType("8910a1ea-4bfe-4321-aa5b-15922b09ad4d", "National ID")
	UNHCRID        = identificationDocumentType("6833cb6d-593f-4f3f-926d-498be74352d1", "UNHCR ID")
	Passport       = identificationDocumentType("567d04e5-abf4-4899-848f-0395264309f0", "Passport")

	BoDiddleyPassport  = identificationDocument("20d194d6-a1ac-483e-8c24-38b5efbaca6f", BoDiddley.ID, "A0JBODIDDLEY129", Passport.ID)
	MaryPoppinsUNHRCID = identificationDocument("0244b59e-5d5c-4e13-af96-da1ccf4e9499", MaryPoppins.ID, "LLP987MARYPOPPINS99", UNHCRID.ID)
	JohnDoeNationalID  = identificationDocument("4c9477c9-c149-4db7-928c-f5e5f915e018", JohnDoe.ID, "B811HJOHNDOE01", NationalID.ID)
)
