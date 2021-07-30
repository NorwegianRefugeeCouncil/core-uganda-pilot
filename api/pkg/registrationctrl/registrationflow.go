package registrationctrl

type RegistrationFlow struct {
	// TODO Country
	TeamID string
	Steps  []Step
}

type RegistrationStepType string

var (
	IndividualAttributes RegistrationStepType = "individual attributes"
	CaseType             RegistrationStepType = "case type"
)

type Step struct {
	Type  RegistrationStepType
	Label string
	Ref   string
}

func (r *RegistrationFlow) toProgress(s state) Progress {
	var progress Progress
	// Handle virgin registration by setting all stages to Unopened
	if s.isVirgin {
		for _, step := range r.Steps {
			stage := Stage{
				Type:        step.Type,
				Ref:         step.Ref,
				StageStatus: Unopened,
			}
			progress = append(progress, stage)
		}
		return progress
	}
	// Otherwise iterate over the steps and evaluate
	for _, step := range r.Steps {
		var stage = Stage{
			Type: step.Type,
			Ref:  step.Ref,
		}
		switch step.Type {
		case IndividualAttributes:
			// If s.isVirgin == false then this step must be done
			stage.StageStatus = Closed
		case CaseType:
			stage.StageStatus = Unopened
			// Examine the given state, if there's a match, update the status.
			for _, closedCase := range s.closedCases {
				if step.Ref == closedCase.CaseTypeID {
					stage.StageStatus = Closed
					break
				}
			}
			for _, opencase := range s.openCases {
				if step.Ref == opencase.CaseTypeID {
					stage.StageStatus = Open
					break
				}
			}
		}
		progress = append(progress, stage)
	}
	return progress
}
