package registrationcontroller

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/seeder"
)

type RegistrationFlow struct {
	TeamID string
	Steps  []*cms.CaseType
}

type RegistrationType string

var ()

type Step struct {
	Type string
}

var (
	Uganda = RegistrationFlow{
		TeamID: "",
		Steps: []*cms.CaseType{
			&seeder.UGIndividualAssessmentCaseType,
			&seeder.UGSituationalAnalysisCaseType,
		},
	}
)
