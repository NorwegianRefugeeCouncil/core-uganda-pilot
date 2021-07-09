package validator

import "context"

type Validate struct {
}

func New() *Validate {
	return new(Validate)
}

func (v *Validate) Struct(ctx context.Context, s interface{}) error {

	return nil
}

type ValidationErrors struct {
	Inputs map[string]Status
}

type Status struct {
	IsValid bool
	Message string
	Causes  []string
}

func (r *Response) String() string {
	if r.Payload == nil {
		return "valid"
	}
	return "invalid"
}
