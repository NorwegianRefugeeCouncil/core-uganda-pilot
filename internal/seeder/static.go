package seeder

import (
	"github.com/nrc-no/core/internal/registrationctrl"
	"github.com/nrc-no/core/pkg/cms"
	"github.com/nrc-no/core/pkg/iam"
)

var (
	individuals                 []iam.Individual
	staffers                    []iam.Staff
	memberships                 []iam.Membership
	countries                   []iam.Country
	nationalities               []iam.Nationality
	relationships               []iam.Relationship
	cases                       []cms.Case
	identificationDocumentTypes []iam.IdentificationDocumentType
	identificationDocuments     []iam.IdentificationDocument

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
)
