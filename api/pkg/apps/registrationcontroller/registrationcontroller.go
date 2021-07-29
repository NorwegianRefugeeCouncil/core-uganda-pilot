package registrationcontroller

import (
	"fmt"
	"github.com/nrc-no/core/pkg/apps/cms"
)

// This package is responsible for interpreting the state of an cms.Individual to determine an Individual's current
// status and what actions a user can take.

// CaseHandler retrieves state information for an individual.
type CaseHandler interface {
	GetAllCases() []*cms.Case
	GetOpenCases() []*cms.Case
	GetClosedCases() []*cms.Case
}

// Controller describes a controller's external interface.
type Controller interface {
	Status() *Status
	Actions() []Action
}

// Action describes an applicable action based on the provided state.
type Action struct {
	Type  string
	Label string
	Link  string
}

// Status indicates the Individuals progress along the CaseFlow.
type Status struct {
	Progress Progress
	Label    string
}

type Progress []Stage

// Stage represents the status of one step along the CaseFlow.
type Stage struct {
	CaseType *cms.CaseType
	CaseStatus
}

// CaseStatus indicates the current status of a Case (either Unopened, Open or Closed)
type CaseStatus int
var (
	Unopened CaseStatus = 0
	Open     CaseStatus = 1
	Closed   CaseStatus = 2
)

type state struct {
	cases       []*cms.Case
	openCases   []*cms.Case
	closedCases []*cms.Case
}

// Ensure RegistrationController implements the Controller interface
var _ Controller = &RegistrationController{}

// RegistrationController compares a CaseHandler against a CaseFlow to provide available Actions and Status of an Individual in the registration process.
type RegistrationController struct {
	handler  CaseHandler
	caseFlow CaseFlow
	state    state
	status   *Status
}

// Implementation

func (r *RegistrationController) Status() *Status {
	status := &Status{}
	progress := r.progress()
	status.Progress = progress
	done := len(r.state.closedCases)
	total := len(r.caseFlow.CaseTypes)
	status.Label = fmt.Sprintf("%d of %d", done, total)
	r.status = status
	return status
}

func (r *RegistrationController) Actions() []Action {
	actions := []Action{}
	for _, p := range r.status.Progress {
		if p.CaseStatus == Open {
			actions = append(actions, Action{
				Type:  p.CaseType.Name,
				Label: "Complete this form",
				Link:  p.CaseType.ID,
			})
		} else if p.CaseStatus == Unopened && len(actions) == 0 {
			// If this case is unopened and there are no open cases, this is the next action
			actions = append(actions, Action{
				Type:  p.CaseType.Name,
				Label: "Continue",
				Link:  p.CaseType.ID,
			})
			// we don't want anything else after this, one step at a time.
			break
		}
	}
	return actions
}

func NewRegistrationController(handler CaseHandler, caseFlow CaseFlow) *RegistrationController {
	controller := &RegistrationController{
		handler:  handler,
		caseFlow: caseFlow,
		state: state{
			cases:       handler.GetAllCases(),
			openCases:   handler.GetOpenCases(),
			closedCases: handler.GetClosedCases(),
		},
	}
	controller.status = controller.Status()
	return controller
}

func (r *RegistrationController) progress() []Stage {
	var progress []Stage
	// Get cases that match the flow
	for _, caseType := range r.caseFlow.CaseTypes {
		for _, closedCase := range r.state.closedCases {
			if caseType.ID == closedCase.CaseTypeID {
				progress = append(progress, Stage{
					CaseType:   caseType,
					CaseStatus: Closed,
				})
				break
			}
		}
		for _, openCase := range r.state.openCases {
			if caseType.ID == openCase.CaseTypeID {
				progress = append(progress, Stage{
					CaseType:   caseType,
					CaseStatus: Open,
				})
				break
			}
		}
		progress = append(progress, Stage{
			CaseType:   caseType,
			CaseStatus: Unopened,
		})
	}
	return progress
}
