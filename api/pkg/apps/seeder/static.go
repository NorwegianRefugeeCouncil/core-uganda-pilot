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
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Does the beneficiary have any restrictions to be referred?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Has restrictions?"}},
						Value: "yes",
					},
				},
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
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Was the referral accepted by the other provider?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Referral accepted"}},
						Value: "yes",
					},
				},
			},
			{
				Name:  "pertinentDetails",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Provide any pertinent details on service needs / requests."}},
			},
		},
	}
	// - Kampala ICLA Team
	UGICLAIndividualIntake = form.Form{
		Controls: []form.Control{
			{
				Name:    "modalityOfService",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Modality of service delivery"}},
				Options: []i18n.Strings{{{"en", "ICLA Legal Aid Centre"}}, {{"en", "Mobile visit"}}, {{"en", "Home visit"}}, {{"en", "Transit Centre"}}, {{"en", "Hotline"}}, {{"en", "Other"}}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:    "livingSituation",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "Living situation"}},
				Options: []i18n.Strings{{{"en", "Lives alone"}}, {{"en", "Lives with family"}}, {{"en", "Hosted by relatives"}}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:        "commentLivingSituation",
				Type:        form.Textarea,
				Label:       i18n.Strings{{"en", "Comment on living situation"}},
				Description: i18n.Strings{{"en", "Additional information, observations, concerns, etc."}},
			},
			{
				Name:    "iclaMeansOfDiscovery",
				Type:    form.Dropdown,
				Label:   i18n.Strings{{"en", "How did you learn about ICLA services?"}},
				Options: []i18n.Strings{{{"en", "ICLA in-person information session"}}, {{"en", "ICLA social media campaign, activities, brochures"}}, {{"en", "ICLA text messages"}}, {{"en", "Another beneficiary/friend/relative"}}, {{"en", "Another organisation"}}, {{"en", "General social media"}}, {{"en", "NRC employee"}}, {{"en", "State authority"}}, {{"en", "Other"}}},
				Validation: form.ControlValidation{
					Required: true,
				},
			},
			{
				Name:        "vulnerability",
				Type:        form.Textarea,
				Label:       i18n.Strings{{"en", "Vulnerability data"}},
				Description: i18n.Strings{{"en", "As needed within a particular context and required for the case"}},
			},
			{
				Name:        "representativeFullName",
				Type:        form.Text,
				Label:       i18n.Strings{{"en", "Full name of representative"}},
				Description: i18n.Strings{{"en", "Lawyer or other person"}},
			},
			{
				Name:        "otherPersonalInfo",
				Type:        form.Textarea,
				Label:       i18n.Strings{{"en", "Other personal information"}},
				Description: i18n.Strings{{"en", "Other personal data as needed to identify the representative within the particular context"}},
			},
			{
				Name:  "reasonForRepresentative",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Reason for representative"}},
			},
			{
				Name:        "guardianshipIsLegal",
				Type:        form.Checkbox,
				Label:       i18n.Strings{{"en", "Is the guardianship legal as per national legislation?"}},
				Description: i18n.Strings{{"en", "If 'yes', attach/upload the legal assessment. If 'no', request or assist in identifying an appropriate legal guardian to represent beneficiary"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Guardianship is legal"}},
					},
				},
			},
			{
				Name:  "capacityToConsent",
				Type:  form.Checkbox,
				Label: i18n.Strings{{"en", "Does the beneficiary have the legal capacity to consent?"}},
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "Beneficiary has legal capacity to consent"}},
					},
				},
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
				Name:        "agreedFollowUp",
				Type:        form.Text,
				Label:       i18n.Strings{{"en", "Agreed follow up with the beneficiary"}},
			},
		}
  }

	// Case Types for Uganda
	// - Kampala Response Team
	UGSituationalAnalysisCaseType      = caseType("0ae90b08-6944-48dc-8f30-5cb325292a8c", "Situational Analysis (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGSituationAnalysis, true)
	UGIndividualResponseCaseType       = caseType("2f909038-0ce4-437b-af17-72fc5d668b49", "Response (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGIndividualResponse, true)
	UGReferralCaseType                 = caseType("ecdaf47f-6fa9-48c8-9d10-6324bf932ed7", "Referral (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGReferral, false)
	UGExternalReferralFollowupCaseType = caseType("2a1b670c-6336-4364-b89d-0e65fc771659", "External Referral Followup (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGExternalReferralFollowup, false)
	UGProtectionIntakeCaseType         = caseType("da20a49d-3cc9-413c-89b8-ff40e3afe95c", "Intake (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionIntake, true)
	UGProtectionFollowUpCaseType       = caseType("dcebe6c8-47cd-4e0f-8562-5680573aed88", "Follow up (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionFollowUp, false)
	// - Kampala ICLA Team
	UGICLAIndividualIntakeCaseType   = caseType("31fb6d03-2374-4bea-9374-48fc10500f81", "ICLA Individual Intake (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLAIndividualIntake, true)
	UGICLACaseAssessmentCaseType     = caseType("bbf820de-8d10-49eb-b8c9-728993ab0b73", "ICLA Case Assessment (UG ICLA)", iam.IndividualPartyType.ID, UgandaICLATeam.ID, UGICLACaseAssessment, false)
	UGProtectionActionReportCaseType = caseType("f4989460-8e76-4d82-aad5-ed2ad3d3d627", "Action Report (UG Protection/Response)", iam.IndividualPartyType.ID, UgandaProtectionTeam.ID, UGProtectionActionReport, false)

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
