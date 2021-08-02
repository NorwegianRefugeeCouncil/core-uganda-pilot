package registrationctrl

import (
	"fmt"
	"github.com/nrc-no/core/pkg/apps/cms"
)

// This package is responsible for interpreting the state of an cms.Individual to determine an Individual's current
// status and what actions a user can take.
// TODO adapt interface to the attachments API rather than cms?
// TODO for now an unregistered individual will return an empty string in the Link field, and Subsequent steps return their CaseTypeID in the Link field.
//

// IndividualHandler retrieves state information for an individual.
type IndividualHandler interface {
	IndividualExists() bool
	GetOpenCases() []*cms.Case
	GetClosedCases() []*cms.Case
}

// Controller describes a controller's external interface.
type Controller interface {
	Status() *Status
	Action() Action
}

// Action describes an applicable action based on the provided state.
type Action struct {
	Type  string
	Label string
	Link  string
}

// Status indicates the Individuals progress along the RegistrationFlow.
type Status struct {
	Progress Progress
	Label    string
}

// Stage represents the status of one step along the RegistrationFlow.
type Stage struct {
	Type RegistrationStepType
	Ref  string
	StageStatus
}

//Progress represents the overall progress of the individual
type Progress []Stage

func (p Progress) nClosed() int {
	count := 0
	for _, stage := range p {
		if stage.StageStatus == Closed {
			count += 1
		}
	}
	return count
}

// StageStatus indicates the current status of a Case (either Unopened, Open or Closed)
type StageStatus int

var (
	Unopened StageStatus = 0
	Open     StageStatus = 1
	Closed   StageStatus = 2
)

type state struct {
	isVirgin    bool
	openCases   []*cms.Case
	closedCases []*cms.Case
}

// Ensure RegistrationController implements the Controller interface
var _ Controller = &RegistrationController{}

// RegistrationController compares a IndividualHandler against a RegistrationFlow to provide available Actions and Status of an Individual in the registration process.
type RegistrationController struct {
	handler          IndividualHandler
	registrationFlow RegistrationFlow
	state            state
	status           *Status
}

// Implementation

func NewRegistrationController(handler IndividualHandler, caseFlow RegistrationFlow) *RegistrationController {
	var state = state{}
	if handler.IndividualExists() {
		state.openCases = handler.GetOpenCases()
		state.closedCases = handler.GetClosedCases()
	} else {
		state.isVirgin = true
	}
	controller := &RegistrationController{
		handler:          handler,
		registrationFlow: caseFlow,
		state:            state,
	}
	controller.Status()
	return controller
}

// Status computes and returns the Individual's registration status according to their state.
func (r *RegistrationController) Status() *Status {
	status := &Status{}
	status.Progress = r.registrationFlow.toProgress(r.state)
	done := status.Progress.nClosed()
	total := len(r.registrationFlow.Steps)
	status.Label = fmt.Sprintf("%d of %d", done, total)
	r.status = status
	return status
}

// Action returns an Action struct representing the next available action on the given Individual
func (r *RegistrationController) Action() Action {
	var action = Action{}
	for _, p := range r.status.Progress {
		if p.StageStatus == Closed {
			continue
		}
		if p.StageStatus == Open {
			action.Label = "Incomplete"
			action.Link = p.Ref
			break
		}
		if p.StageStatus == Unopened {
			action.Label = "Next step"
			action.Link = p.Ref
			break
		}
	}
	return action
}
