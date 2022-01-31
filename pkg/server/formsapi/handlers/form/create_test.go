package form

import (
	"context"
	"testing"

	"github.com/nrc-no/core/pkg/api/types"
)

func Test_validateRecipientFormParent(t *testing.T) {

	aNonRecipientForm := types.FormDefinition{
		Type: types.DefaultFormType,
	}

	aRecipientFormWithoutParent := types.FormDefinition{
		Type: types.RecipientFormType,
	}

	aRecipientFormWithParent := types.FormDefinition{
		Type: types.RecipientFormType,
		Fields: []*types.FieldDefinition{
			{
				Key: true,
				FieldType: types.FieldType{
					Reference: &types.FieldTypeReference{},
				},
			},
		},
	}

	tests := []struct {
		name    string
		ctx     context.Context
		form    types.FormDefinition
		getForm func(ctx context.Context, formID string) (*types.FormDefinition, error)
		wantErr bool
	}{
		{
			name:    "not a recipient form",
			ctx:     context.Background(),
			form:    aNonRecipientForm,
			getForm: nil,
			wantErr: false,
		}, {
			name:    "no parent",
			ctx:     context.Background(),
			form:    aRecipientFormWithoutParent,
			getForm: nil,
			wantErr: false,
		}, {
			name: "with recipient parent",
			ctx:  context.Background(),
			form: aRecipientFormWithParent,
			getForm: func(ctx context.Context, formID string) (*types.FormDefinition, error) {
				return &types.FormDefinition{Type: types.RecipientFormType}, nil
			},
			wantErr: false,
		}, {
			name: "with non-recipient parent",
			ctx:  context.Background(),
			form: aRecipientFormWithParent,
			getForm: func(ctx context.Context, formID string) (*types.FormDefinition, error) {
				return &types.FormDefinition{Type: types.DefaultFormType}, nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRecipientFormParent(tt.ctx, tt.form, tt.getForm); (err != nil) != tt.wantErr {
				t.Errorf("validateRecipientFormParent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
