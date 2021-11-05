package seeder

import (
	"context"
	"github.com/nrc-no/core/pkg/store"

	"github.com/nrc-no/core/pkg/api/types"
)

var noOpts = []string{}
var yesNoOpts = []string{"Yes", "No"}
var wgShortSet = []string{"Moderate Impairment", "Severe Impairment"}
var globalDisplacementStatuses = []string{"Refugee", "Internally Displaced", "Host Community", "Other"}
var ugGenders = []string{"Male", "Female"}
var ugServices = []string{
	"Health care (including medication)",
	"Legal assistance",
	"Education",
	"Mental Health",
	"Transportation",
	"Food",
	"Non-food items (including hygiene items)",
	"Disability",
	"MPC",
	"Shelter/Housing",
	"Shelter construction/repair",
	"Youth Livelihoods (e.g. vocational training)",
	"Small/Medium Business Grants",
	"Other livelihood activities",
}
var ugICLALegalIssues = []string{
	"RSD",
	"ELP",
	"HLP",
	"IDP registration",
	"Other",
}
var ugPriority = []string{
	"High",
	"Medium",
	"Low",
}
var ugICLADisplacementStatus = []string{
	"Unregistered asylum seeker",
	"Registered asylum seeker",
	"Refugee",
}
var ugICLARSDDocuments = []string{
	"Family Attesttion",
	"Refugee ID",
	"Asylum certificate",
	"Rejection decision",
	"Other",
}
var ugICLASpecificHLPNeed = []string{
	"Housing",
	"Land",
	"Property",
}
var ugICLAHomeOwnership = []string{
	"Own house",
	"Rent",
	"Other",
}
var ugICLAEvictionDocuments = []string{
	"Eviction Notice",
	"Other",
}
var ugICLANatureLandTenure = []string{
	"Joint ownership",
	"Co-ownership",
	"Individual ownership",
	"Other",
}
var ugICLANatureTenure = []string{
	"Mailo",
	"Lease",
	"Freehold",
	"Sustomary",
}
var ugICLATypeOfDocumentation = []string{
	"Legal",
	"Civil",
}
var ugICLATypeOfAgreement = []string{
	"Oral",
	"Written",
}
var ugICLATypeOfChallenge = []string{
	"Employment",
	"Business",
}
var ugICLATypeOfActions = []string{
	"(1) Discussion with supervisor/team leader",
	"(2) Conducting legal analysis, including the study of judicial practice",
	"(3) Preparing letters, inquiries to various authorities",
	"(4) Drafting of other legal documents (such leases or contracts)",
	"(5) Lodging of a court application",
	"(6) Attending of court session/hearing",
	"(7) Review of the decision/appeal",
	"(8) Negotiation",
	"(9) Follow up with relevant administrative authority or other entities",
	"(10) Accompaniment",
	"(11) Other",
}
var ugICLATypesOfServices = []string{
	"Legal counselling",
	"Referral",
	"Relocation",
	"Livelihood",
	"Business support",
}
var ugICLACaseClosureReason = []string{
	"All action plan objectives achieved",
	"Other priorities by beneficiary",
	"Beneficiary unreachable",
	"Other",
}
var ugICLAAgreedFollowupMeans = []string{
	"Schedule in-person meeting with beneficiary",
	"Schedule phone call",
	"Other",
}

func buildField(name, description string, options []string, key, required bool, fieldType types.FieldType) *types.FieldDefinition {
	return &types.FieldDefinition{
		Name:        name,
		Description: description,
		Key:         key,
		Required:    required,
		FieldType:   fieldType,
	}
}

func buildForm(databaseId, folderId, name string, fields []*types.FieldDefinition) *types.FormDefinition {
	return &types.FormDefinition{
		DatabaseID: databaseId,
		FolderID:   folderId,
		Name:       name,
		Fields:     fields,
	}
}

func seedUgandaForms(ctx context.Context, formStore store.FormStore, ugandaDatabaseId, coFolderId, iclaFolderId, intakeFolderId, protectionFolderId string) error {
	// Root entities ---------------------------------

	rootIndividualForm, err := formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		coFolderId,
		"Individual",
		[]*types.FieldDefinition{
			buildField("Full Name", "The full name of the individual", noOpts, false, true, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
			buildField("Preferred Name", "The name which will be used to refer to the beneficiary within Core", noOpts, false, true, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
		},
	))
	if err != nil {
		return err
	}

	rootHouseholdForm, err := formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		coFolderId,
		"Household",
		[]*types.FieldDefinition{
			buildField("Household Name", "", noOpts, false, true, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
		},
	))
	if err != nil {
		return err
	}

	// Global Intake ---------------------------------

	individualBeneficiaryForm, err := formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		coFolderId,
		"Individual Beneficiary",
		[]*types.FieldDefinition{
			buildField("Individual", "Individual who is being registered as a beneficiary", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     rootIndividualForm.ID,
				},
			}),
			buildField("Consent", "Did the beneficiary consent to NRC using their data", yesNoOpts, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Consent URL", "Link to proof of consent", noOpts, false, true, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
			buildField("Anonymous", "Did the beneficiary request to remain anonymous", yesNoOpts, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Minor", "Is this beneficiary a minor", yesNoOpts, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Protection Concern", "Does this beneficiary present protection concerns", yesNoOpts, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Physical Disability", "", yesNoOpts, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Physical Disability - Explanation", "", wgShortSet, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Sensory Disability", "", yesNoOpts, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Sensory Disability - Explanation", "", wgShortSet, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Mental Disability", "", yesNoOpts, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Mental Disability - Explanation", "", wgShortSet, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Displacement Status", "", globalDisplacementStatuses, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Gender", "", ugGenders, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Household", "Household to which this beneficiary belongs", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     rootHouseholdForm.ID,
				},
			}),
			buildField("Head of Household", "Is this beneficiary the head of their household", yesNoOpts, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
		},
	))
	if err != nil {
		return err
	}

	// UG Intake ---------------------------------
	var protectionPriorityOpts = []string{"High (follow-up requested in 24 hours)", "Medium (follow-up in 3 days)", "Low (follow-up in 7 days)"}
	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		intakeFolderId,
		"Uganda Individual Intake Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this intake form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Situation Analysis", "A form used to carry out a situation analysis for an individual in Uganda", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Safe Dignified Life", "Do you think you are living a safe and dignified life? Are you achieving what you want? Are you able to live a good life? Probe for description", noOpts, false, true, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("Addressing Challenges/Barriers", "How are you addressing these challenges and barriers? What is standing in your way? Can you give me some examples of how you are dealing with these challenges?", noOpts, false, true, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("Possible Solutions", "What are some solutions you see for this and how could we work together on these solutions? How could we work to reduce these challenges together?", noOpts, false, true, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("What Could We Do Together", "If we were to work together on this, what could we do together? What would make the most difference for you?", noOpts, false, true, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
					},
				},
			}),
			buildField("Response", "Response form for an individual in Uganda", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Which service has the individual/community requested as a starting point of support?", "", ugServices, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("Provide any pertinent details on service needs / requests", "", noOpts, false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What other services has the individual /household requested/identified?", "", ugServices, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("Provide any pertinent details on service needs / requests for the other services requested.", "", noOpts, false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What is the perceived priority response level of the individual / household?", "", ugPriority, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("Provide any pertinent details on how priority was determined.", "", noOpts, false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Other information", "Comments or notes", noOpts, false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
					},
				},
			}),
		},
	))
	if err != nil {
		return err
	}

	// UG Protection ---------------------------------

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		protectionFolderId,
		"Uganda Protection Case Opening Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this protection intake form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Protection Intake", "A form used to collect intake details for the Uganda Protection team", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Date of screening", "", noOpts, false, true, types.FieldType{
							Date: &types.FieldTypeDate{},
						}),
						buildField("Have you been exposed to any protection risks?", "", yesNoOpts, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("What type of protection concern experienced?", "", []string{"Physical violence", "Neglect", "Family separation", "Arrest", "Denial of resources", "Psychological violence"}, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("Provide details", "Narrative", noOpts, false, true, types.FieldType{
							MultilineText: &types.FieldTypeMultilineText{},
						}),
						buildField("Response priority", "", protectionPriorityOpts, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
					},
				},
			}),
			buildField("Social Status Assessment", "Used to collect information regarding the social status of an individual beneficiary", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Main obstacles you face in meeting accomodation needs", "", []string{"Insufficient funds", "Distance issues", "Insecurity", "Social discrimination", "Insufficient quantity of goods", "Inadequate quality of goods/services", "Insufficient capabilities and competences", "Other"}, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("Can the HH meet WASH needs?", "", []string{"We can meet needs without worry", "We can meet our needs", "We can barely meet our needs", "We are unable to meet our needs", "We are totally unable to meet our needs"}, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("Main obstacles in meeting WASH needs", "", []string{"Insufficient funds", "Distance constraints", "Insecurity", "Social discrimination", "Insufficient quantity of goods/services", "Inadequate quality of goods/services", "Other"}, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("Summary narrative", "", noOpts, false, false, types.FieldType{
							MultilineText: &types.FieldTypeMultilineText{},
						}),
					},
				},
			}),
		},
	))
	if err != nil {
		return err
	}

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		protectionFolderId,
		"Uganda Protection Incident Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this incident form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Location of incident", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Time of incident", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Date of incident", "", noOpts, false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
			buildField("Received by", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Vulnerability", "", []string{"Child at risk", "Elder at risk", "Single parent", "Separated child", "Disability", "Woman at risk", "Legal and physical protection", "Medical condition", "Pregnant or lactating woman"}, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("Description of the incident", "i.e. Where, when, what, who involved", noOpts, false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			buildField("Has the incident been reported to the police?", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Comment", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			buildField("Has the incident been reported to", "", []string{"UNCHR", "Other platform"}, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
		},
	))
	if err != nil {
		return err
	}

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		protectionFolderId,
		"Uganda Protection Action Report Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this action report form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Which service has the benficiary together with staff agreed to take?", "", []string{"Cash support", "Referral", "Relocation", "Livelihood", "Business support", "Other"}, false, true, types.FieldType{
				SingleSelect: &types.FieldTypeSingleSelect{},
			}),
			buildField("If \"Other\", specify", "", noOpts, false, false, types.FieldType{
				Text: &types.FieldTypeText{},
			}),
			buildField("Narrate/Describe the response action agreed upon", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			buildField("Agreed follow-up with beneficiary", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
		}))
	if err != nil {
		return err
	}

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		protectionFolderId,
		"Uganda Protection Followup Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this followup form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Follow up after", "", []string{"1 week", "2 weeks", "1 month", "3 months"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Agreed follow-up with the beneficiary", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
		},
	))
	if err != nil {
		return err
	}

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		protectionFolderId,
		"Uganda Protection Referral Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this referral form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Priority", "", protectionPriorityOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Referred via", "", []string{"Phone (High priority only)", "Email", "In person"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Referral date", "", noOpts, false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
			buildField("Receiving agency", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Name of partner case worker", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Position of person receiving referral", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Contact of person receiving referral", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Consent to release information", "", yesNoOpts, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Has the person expressed any restrictions on referrals?", "", noOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("If yes, specify", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			buildField("Is the beneficiary a minor", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Name of the primary caregiver", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Relationship to the child", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Has the caregiver been informed of the referral?", "", yesNoOpts, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("If not, explain", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			buildField("Reason for referral", "", noOpts, false, true, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			buildField("Type of referral", "", []string{"Health", "Livelihood/IGAS", "Pschosocial support", "Safety & security", "Education", "Shelter"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
		},
	))
	if err != nil {
		return err
	}

	// UG ICLA ---------------------------------
	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		iclaFolderId,
		"Uganda Case Opening Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this ICLA intake form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Intake", "A form used to collect intake details for the Uganda ICLA team", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("What legal issue/concern are you facing?", "Describe the concern if possible", noOpts, false, true, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("Select the legal issue of concern", "", ugICLALegalIssues, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("If other, Please specify", "", noOpts, false, false, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("What action has been taken to solve the problem if any", "", noOpts, false, true, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("Is there a representative for this individual", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Full name of representative", "", noOpts, false, false, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("Reason for representative", "", noOpts, false, false, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("Is the guardianship legal as per national legislation", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Attach/upload the legal/court order", "", noOpts, false, false, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("What is the individual's displacement status", "", ugICLADisplacementStatus, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Are you at risk of being stateless", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Describe this in detail", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What RSD documents do you have", "", ugICLARSDDocuments, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Comment on RSD documents", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Specific RSD issues presented", "Narrative", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What specific HLP concern is presented", "", ugICLASpecificHLPNeed, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Does the individual  stay in their own house or rent", "", ugICLAHomeOwnership, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("In case of rent, does the individual posses any  agreement", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("What kind of agreement or  proof does the individual possess", "", noOpts, false, false, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
						buildField("Have you been or are you at risk of eviction", "Housing", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("If Yes, What eviction document or proof do you posses?", "Housing", ugICLAEvictionDocuments, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Comment on eviction document", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Are you the legal owner of the land", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Nature of tenancy", "", ugICLANatureLandTenure, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Nature of tenure", "", ugICLANatureTenure, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Land supporting documents possessed", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Have you been or are you at risk of eviction", "Land", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("If Yes, What eviction document or proof do you posses?", "Land", ugICLAEvictionDocuments, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Specific land issues", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Nature of property in contest", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Do you have legal ownership of property", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Proof of property ownership", "(Supporting documents)", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Inquiry on property acquisition", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What documentation challenges do you have", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Type of document", "", ugICLATypeOfDocumentation, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("What action had been taken so far on this issue", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Is it an employment or business challenge", "", ugICLATypeOfChallenge, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("What employment challenges do you have", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What type of agreement do you have", "", ugICLATypeOfAgreement, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("What actions have been taken", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What business related challenge do you have", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What business registration services do you need", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("What actions have been taken", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
					},
				},
			}),
			buildField("Consent", "Used to collect information regarding the consent of an individual beneficiary", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Consent", "Did the beneficiary consent to NRC ICLA using their data", yesNoOpts, false, true, types.FieldType{
							SingleSelect: &types.FieldTypeSingleSelect{},
						}),
						buildField("Consent URL", "Link to proof of consent", noOpts, false, true, types.FieldType{
							Text: &types.FieldTypeText{},
						}),
					},
				},
			}),
			buildField("Case Assessment - Case Plan", "ICLA case plan for an individual beneficiary", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Type of actions for case worker agreed upon with beneficiary ", "", ugICLATypeOfActions, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Notes on types of actions agreed upon", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Are there any elements of risks for the safety or well-being of the beneficiary or that of a relative in relation to the suggested course of action", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Narrative regarding elements of risks", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Any particular Protection Risks", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Is the guardianship legal as per national legislation", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("If yes, please indicate what type", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Guardianship notes", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Are there any unintended negative consequences of the suggested course of actions for the beneficiary's family or larger community", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Unintended consequences notes", "If any of the answers were 'yes', discuss with the beneficiary what might be done to avoid or minimise the risks or negative consequences", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Actions agreed upon with the beneficiary", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Does the beneficiary agree to continue with the case", "Discuss the pro's and con's of the suggested course of action, including the analysis of risks", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Notes on pro's and con's of suggested course of action", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Is the guardianship legal as per national legislation", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Is a Best Interest Determination needed for the case", "If 'yes', refer the case to social services or an appropriate child protection actor", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Agreed follow up means", "", ugICLAAgreedFollowupMeans, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
					},
				},
			}),
			buildField("Case Assessment - Action Plan", "ICLA action plan for an individual beneficiary", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Which service has the beneficiary together with staff agreed to take", "", ugICLATypesOfServices, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("If other specify", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
						buildField("Type of actions for case worker agreed upon with beneficiary ", "", ugICLATypeOfActions, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Action comment", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
					},
				},
			}),
			buildField("Case Assessment - Case Closure", "Used to close an ICLA case in Uganda", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Case closure date", "", noOpts, false, false, types.FieldType{Date: &types.FieldTypeDate{}}),
						buildField("What is the reason for the case closure?", "", ugICLACaseClosureReason, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Notes", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
					},
				},
			}),
		},
	))
	if err != nil {
		return err
	}

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		iclaFolderId,
		"Uganda ICLA Followup Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this followup form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Date of follow-up", "", noOpts, false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
			buildField("Follow-up", "", []string{"In-person interview", "Phone conversation", "Discussion with supervisor", "Conducting legal analysis", "Preparing letters, inquiries to various authorities", "Drafting of other legal documents (e.g. contracts)", "Lodging of a court application", "Attending of court session", "Review of the decision", "Execution of the court decision", "Negotiation", "Follow-up with relevant administrative authority", "Accompaniement", "Other"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("If \"Other\", specify", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Notes from the follow-up undertaken", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			buildField("Copies of documents", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
		},
	))
	if err != nil {
		return err
	}

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		iclaFolderId,
		"Uganda ICLA Appointment Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this appointment form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Name", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Place", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
			buildField("Date", "", noOpts, false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
			buildField("Preferred contact method", "", []string{"Email", "Telephone", "Other"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Appointment purpose", "", []string{"HLP", "LCD", "RSD", "Employment/Business", "Other"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Preferred date", "", noOpts, false, true, types.FieldType{Date: &types.FieldTypeDate{}}),
		},
	))
	if err != nil {
		return err
	}

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		iclaFolderId,
		"Uganda ICLA Referral Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this ICLA referral form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Beneficiary's information", "", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Type of referral (internal)", "", []string{"Shelter/NFI", "Livelihood/Food", "Security", "Education", "WASH", "Camp management/UDOC"}, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Type of referral (external)", "If yes, provide details below", yesNoOpts, false, false, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Organization", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						buildField("Contact person", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						buildField("Phone number", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						buildField("Email", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						buildField("Types of services/assistence requested", "", []string{"Healthcare (including medication)", "Legal assistance", "Education", "Mental health", "Transportation", "Food", "Non-food items (including hygiene items)", "Disability", "MPC", "Shelter/Housing", "Shelter construction/repair", "Youth livelihood (e.g. vocational training)", "Small/medium business grants", "Other livelihood activities"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("If \"Other\", specify", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						buildField("Reason for the referral", "", noOpts, false, true, types.FieldType{Text: &types.FieldTypeText{}}),
					},
				},
			}),
			buildField("Consent", "", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Has the beneficiary consented to the release of their information for the referral?", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("If 'yes', upload a signed consent form and proceed", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
						buildField("If 'no', explain the reason why and do not refer the case", "", noOpts, false, false, types.FieldType{Text: &types.FieldTypeText{}}),
					},
				},
			}),
			buildField("Means of referral", "", noOpts, false, true, types.FieldType{
				SubForm: &types.FieldTypeSubForm{
					Fields: []*types.FieldDefinition{
						buildField("Does the beneficiary have any restrictions to being referred?", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Means of referral", "", []string{"Phone", "Email", "Personal meeting", "Other"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
						buildField("Means and terms of receiving feedback from the client", "", []string{"Phone", "Email", "Personal meeting", "Other"}, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
					},
				},
			}),
		},
	))
	if err != nil {
		return err
	}

	// // UG Other ---------------------------------

	_, err = formStore.Create(ctx, buildForm(
		ugandaDatabaseId,
		coFolderId,
		"Uganda External Referral Form",
		[]*types.FieldDefinition{
			buildField("Individual Beneficiary", "The beneficiary this external referral form has been completed for", noOpts, false, true, types.FieldType{
				Reference: &types.FieldTypeReference{
					DatabaseID: ugandaDatabaseId,
					FormID:     individualBeneficiaryForm.ID,
				},
			}),
			buildField("Was the referral accepted by the other provider", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
			buildField("Provide any pertinent details on service needs / requests", "", noOpts, false, false, types.FieldType{MultilineText: &types.FieldTypeMultilineText{}}),
			buildField("This case is now closed", "", yesNoOpts, false, true, types.FieldType{SingleSelect: &types.FieldTypeSingleSelect{}}),
		},
	))
	if err != nil {
		return err
	}

	return nil
}
