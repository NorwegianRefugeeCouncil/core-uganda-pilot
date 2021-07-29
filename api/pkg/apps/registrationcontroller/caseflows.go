package registrationcontroller

import (
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/seeder"
)

type CaseFlow struct {
	TeamID    string
	CaseTypes []*cms.CaseType
}

var (
	Uganda = CaseFlow{
		TeamID: "",
		CaseTypes: []*cms.CaseType{
			&seeder.UGIndividualAssessmentCaseType,
			&seeder.UGSituationalAnalysisCaseType,
			&seeder.UGReferralCaseType,
			&seeder.UGExternalReferralFollowupCaseType,
		},
	}
)
