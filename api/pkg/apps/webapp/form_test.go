package webapp

import (
	"github.com/nrc-no/core/pkg/validation"
	"html/template"
	"testing"
)

func TestForm_RenderValidationError(t *testing.T) {
	type fields struct {
		WasValidated     bool
		ValidationErrors validation.ErrorList
		Fields           []FormField
	}
	type args struct {
		field string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   template.HTML
	}{
		{
			name: "single error",
			fields: fields{
				WasValidated: true,
				ValidationErrors: validation.ErrorList{&validation.Error{
					Type:     validation.ErrorTypeInvalid,
					Field:    ".mock",
					BadValue: "%",
					Detail:   "mock",
				}},
				Fields: nil,
			},
			args: args{"mock"},
			want: template.HTML(`<div id="mockFeedback" class="invalid-feedback"><p>mock</p></div>`),
		},
		{
			name: "multi error",
			fields: fields{
				WasValidated: true,
				ValidationErrors: validation.ErrorList{&validation.Error{
					Type:     validation.ErrorTypeInvalid,
					Field:    ".mock",
					BadValue: "",
					Detail:   "mock",
				}, &validation.Error{
					Type:     validation.ErrorTypeInvalid,
					Field:    ".mock",
					BadValue: nil,
					Detail:   "mock2",
				}},
				Fields: nil,
			},
			args: args{"mock"},
			want: template.HTML(`<div id="mockFeedback" class="invalid-feedback"><p>mock</p><p>mock2</p></div>`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Form{
				WasValidated:     tt.fields.WasValidated,
				ValidationErrors: tt.fields.ValidationErrors,
				Fields:           tt.fields.Fields,
			}
			if got := f.RenderValidationError(tt.args.field); got != tt.want {
				t.Errorf("RenderValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}
