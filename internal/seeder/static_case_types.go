package seeder

import (
	"github.com/nrc-no/core/internal/form"
	"github.com/nrc-no/core/internal/i18n"
	"github.com/nrc-no/core/pkg/cms"
	"github.com/nrc-no/core/pkg/iam"
)

var (
	caseTypes []cms.CaseType

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

	UGICLACaseClosure = form.Form{
		Controls: []form.Control{
			{
				Name:  "closureDate",
				Type:  form.Date,
				Label: i18n.Strings{{"en", "Case closure date"}},
			},
			{
				Name:     "reasonForClosure",
				Type:     form.Dropdown,
				Label:    i18n.Strings{{"en", "What is the reason for the case closure?"}},
				Multiple: true,
				Options: []i18n.Strings{
					{{"en", "All action plan objectives achieved"}},
					{{"en", "Other priorities by beneficiary"}},
					{{"en", "Beneficiary unreachable"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "closureNotes",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Notes"}},
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
	UGICLAIndividualCasePlan = form.Form{
		Controls: form.Controls{
			{
				Name:     "actionTypes",
				Type:     form.Dropdown,
				Multiple: true,
				Label:    i18n.Strings{{"en", "Type of actions for case worker agreed upon with beneficiary"}},
				Options: []i18n.Strings{
					{{"en", "Discussion with supervisor/team leader"}},
					{{"en", "Conducting legal analysis, including the study of judicial practice"}},
					{{"en", "Preparing letters/inquiries to various authorities"}},
					{{"en", "Drafting of other legal documents (such as leases or contracts"}},
					{{"en", "Lodging of a court application"}},
					{{"en", "Attending of court session/hearing"}},
					{{"en", "Review of the decision/appeal"}},
					{{"en", "Negotiation"}},
					{{"en", "Follow-up with relevant administrative authority or other entity"}},
					{{"en", "Accompaniment"}},
					{{"en", "Other"}},
				},
			},
			{
				Name:  "notes",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Notes/Narrative"}},
			},
			{
				Type:  form.Subtitle,
				Label: i18n.Strings{{"en", "Individual Risk Assessment"}},
			},
			{
				Name:  "riskElementsExist",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Are there any elements of risk for the safety or well-being of the beneficiary or that of a relative in relation to the suggested course of action?"}},
			},
			{
				Name:  "riskNarrative",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Notes/Narrative"}},
			},
			{
				Name:  "particularProtectionRisks",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Any particular Protection Risks?"}},
			},
			{
				Name:  "typeOfProtectionRisks",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "If yes, indicate what type"}},
			},
			{
				Name:  "protectionRisksNarrative",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Narrative"}},
			},
			{
				Name:  "hasNegativeConsequences",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Are there any unintended consequences of the suggested course of actions for the beneficiary's family or larger community?"}},
			},
			{
				Name:  "consequencesNarrative",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Notes/Narrative"}},
			},
			{
				Type:  form.Hint,
				Label: i18n.Strings{{"en", "If any of the answers were 'yes', discuss with the beneficiary what might be done to avoid or minimise the risks or negative consequences."}},
			},
			{
				Name:  "agreedUponActions",
				Type:  form.Text,
				Label: i18n.Strings{{"en", "Actions agreed upon with the beneficiary"}},
			},
			{
				Name:  "actionsProsCons",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Discuss the pro's and con's of the suggested course of action, including the analysis of risks. Does the beneficiary agree to continue with the case?"}},
			},
			{
				Name:  "prosConsNarrative",
				Type:  form.Textarea,
				Label: i18n.Strings{{"en", "Notes/Narrative"}},
			},
			{
				Type:  form.Subtitle,
				Label: i18n.Strings{{"en", "Best Interest Determination"}},
			},
			{
				Type:  form.Hint,
				Label: i18n.Strings{{"en", "In case of a minor, or person with limited (mental) capacity"}},
			},
			{
				Name:  "BIDNeeded",
				Type:  form.Boolean,
				Label: i18n.Strings{{"en", "Is a Best Interest Determination needed for the case?"}},
			},
			{
				Type:  form.Hint,
				Label: i18n.Strings{{"en", "If yes, refer the case to social services or an appropriate child protection actor"}},
			},
			{
				Name: "followUpMeans",
				Type: form.Checkbox,
				CheckboxOptions: []form.CheckboxOption{
					{
						Label: i18n.Strings{{"en", "In-person meeting"}},
						Value: "in person",
					},
					{
						Label: i18n.Strings{{"en", "Phone call"}},
						Value: "phone call",
					},
					{
						Label: i18n.Strings{{"en", "Other"}},
						Value: "other",
					},
				},
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
	UGSituationalAnalysisCaseType = caseType("0ae90b08-6944-48dc-8f30-5cb325292a8c", "🌎 Registration - Situational Analysis", iam.IndividualPartyType.ID, KampalaRegistrationTeam.ID, UGSituationAnalysis, true)
	UGIndividualResponseCaseType  = caseType("2f909038-0ce4-437b-af17-72fc5d668b49", "🌎 Registration - Response", iam.IndividualPartyType.ID, KampalaRegistrationTeam.ID, UGIndividualResponse, true)

	UGProtectionIntakeCaseType                 = caseType("da20a49d-3cc9-413c-89b8-ff40e3afe95c", "🤲 Protection - Intake", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionIntake, true)
	UGProtectionFollowUpCaseType               = caseType("dcebe6c8-47cd-4e0f-8562-5680573aed88", "🤲 Protection - Follow up", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionFollowUp, false)
	UGProtectionSocialStatusAssessmentCaseType = caseType("e3b30f91-7181-41a3-8187-f176084a0ab2", "🤲 Protection - Social Status Assessment", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionSocialStatusAssessment, false)
	UGProtectionReferralCaseType               = caseType("dc18bf9d-e812-43a8-b843-604c23306cd6", "🤲 Protection - UG Protection Referral", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionReferral, false)
	UGProtectionIncidentCaseType               = caseType("f6117a29-db5a-49d7-b564-bf42740ae824", "🤲 Protection - Incident", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionIncident, false)
	UGProtectionActionReportCaseType           = caseType("f4989460-8e76-4d82-aad5-ed2ad3d3d627", "🤲 Protection - Action Report", iam.IndividualPartyType.ID, KampalaProtectionTeam.ID, UGProtectionActionReport, false)

	// - Kampala ICLA Team
	UGICLAFollowUpCaseType           = caseType("415be6d4-cf1b-484a-9bad-83acd8474498", "⚖️ ICLA - Follow up", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAFollowUp, false)
	UGICLAIntakeCaseType             = caseType("61fb6d03-2374-4bea-9374-48fc10500f81", "⚖️ ICLA - Intake", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAIntake, true)
	UGICLACaseAssessmentCaseType     = caseType("bbf820de-8d10-49eb-b8c9-728993ab0b73", "⚖️ ICLA - Case Assessment", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLACaseAssessment, false)
	UGICLAAppointmentCaseType        = caseType("27064ded-fbfe-4197-830c-164a797d5306", "⚖️ ICLA - Appointment", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAAppointment, false)
	UGICLAConsentCaseType            = caseType("3ad2d524-4dd0-4834-9fc2-47808cf66941", "⚖️ ICLA - Consent", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAConsent, false)
	UGICLAActionPlanCaseType         = caseType("2b4f46a7-aebd-4754-89fd-dc7897a79ddb", "⚖️ ICLA - Action Plan", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAActionPlan, false)
	UGICLACaseClosureCaseType        = caseType("2411793d-55a0-46af-b4f6-a2310d66568f", "⚖️ ICLA - Case Closure", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLACaseClosure, false)
	UGICLAReferralCaseType           = caseType("9896c0f1-8d66-4657-92f2-e67a7afcf9ab", "⚖️ ICLA - Referral", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAReferral, false)
	UGICLAIndividualCasePlanCaseType = caseType("80b2b596-1664-47ff-975b-b4c5af23abdf", "⚖️ ICLA - Individual Case Plan & Risk Assessment", iam.IndividualPartyType.ID, KampalaICLATeam.ID, UGICLAIndividualCasePlan, false)
)
