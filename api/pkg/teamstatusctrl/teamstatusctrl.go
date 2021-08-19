package teamstatusctrl

import (
	"github.com/nrc-no/core/pkg/apps/cms"
)

type TeamStatusCtrl struct {
	CasesForIndividual  []*cms.Case
	TeamIntakeCaseTypes []*cms.CaseType
}

type TeamStatusAction struct {
	label      string
	caseTypeId string
	redirectId string
}

func (tsc *TeamStatusCtrl) GetTeamStatusActions() []TeamStatusAction {
	teamStatusActions := []TeamStatusAction{}

	// Create actions for team intake case type cases that
	// already exist for the individual
	for _, kase := range tsc.CasesForIndividual {
		for _, ct := range tsc.TeamIntakeCaseTypes {
			if kase.CaseTypeID == ct.ID {
				newAction := TeamStatusAction{
					label:      ct.Name,
					caseTypeId: ct.ID,
					redirectId: kase.ID,
				}
				teamStatusActions = append(teamStatusActions, newAction)
				break
			}
		}
	}

	// Create actions for team intake case types not yet
	// accounted for in the individual's case list
	for _, ct := range tsc.TeamIntakeCaseTypes {
		ctAccountedForInActions := false

		for _, tsa := range teamStatusActions {
			if tsa.caseTypeId == ct.ID {
				ctAccountedForInActions = true
				break
			}
		}

		if !ctAccountedForInActions {
			newAction := TeamStatusAction{
				label:      ct.Name,
				caseTypeId: ct.ID,
				redirectId: "",
			}

			teamStatusActions = append(teamStatusActions, newAction)
		}
	}

	return teamStatusActions
}
