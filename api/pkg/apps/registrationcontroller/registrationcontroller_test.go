package registrationcontroller

import (
	"reflect"
	"testing"
)

func TestNewRegistrationController(t *testing.T) {
	type args struct {
		handler  CaseHandler
		caseFlow CaseFlow
	}
	tests := []struct {
		name string
		args args
		want *RegistrationController
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRegistrationController(tt.args.handler, tt.args.caseFlow); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRegistrationController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistrationController_Actions(t *testing.T) {
	type fields struct {
		handler  CaseHandler
		caseFlow CaseFlow
		state    state
		status   *Status
	}
	tests := []struct {
		name   string
		fields fields
		want   []Action
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RegistrationController{
				handler:  tt.fields.handler,
				caseFlow: tt.fields.caseFlow,
				state:    tt.fields.state,
				status:   tt.fields.status,
			}
			if got := r.Actions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistrationController_Status(t *testing.T) {
	type fields struct {
		handler  CaseHandler
		caseFlow CaseFlow
		state    state
		status   *Status
	}
	tests := []struct {
		name   string
		fields fields
		want   *Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RegistrationController{
				handler:  tt.fields.handler,
				caseFlow: tt.fields.caseFlow,
				state:    tt.fields.state,
				status:   tt.fields.status,
			}
			if got := r.Status(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistrationController_progress(t *testing.T) {
	type fields struct {
		handler  CaseHandler
		caseFlow CaseFlow
		state    state
		status   *Status
	}
	tests := []struct {
		name   string
		fields fields
		want   []Stage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RegistrationController{
				handler:  tt.fields.handler,
				caseFlow: tt.fields.caseFlow,
				state:    tt.fields.state,
				status:   tt.fields.status,
			}
			if got := r.progress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("progress() = %v, want %v", got, tt.want)
			}
		})
	}
}
