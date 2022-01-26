package store

import (
	"testing"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/testutils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

// TestFlattenHydrateFormDefinition tests that we can flatten and re-hydrate a FormDefinition
// with various field types
func TestFlattenHydrateFormDefinition(t *testing.T) {

	aFormWithFields := func(fields ...*types.FieldDefinition) *types.FormDefinition {
		return &types.FormDefinition{
			ID:         "formId",
			DatabaseID: "databaseId",
			FolderID:   "folderId",
			Name:       "formName",
			Fields:     fields,
		}
	}

	const fieldName = "fieldName"
	const fieldId = "fieldId"

	aFormRef := &types.FormRef{
		DatabaseID: uuid.NewV4().String(),
		FormID:     uuid.NewV4().String(),
	}

	tests := []struct {
		name           string
		formDefinition *types.FormDefinition
	}{
		{
			name: "with text field",
			formDefinition: aFormWithFields(
				testutils.ATextField(
					testutils.FieldID(fieldId),
					testutils.FieldName(fieldName),
				),
			),
		}, {
			name: "with multiline field",
			formDefinition: aFormWithFields(
				testutils.AMultilineTextField(
					testutils.FieldID(fieldId),
					testutils.FieldName(fieldName),
				),
			),
		}, {
			name: "with date field",
			formDefinition: aFormWithFields(
				testutils.ADateField(
					testutils.FieldID(fieldId),
					testutils.FieldName(fieldName),
				),
			),
		}, {
			name: "with month field",
			formDefinition: aFormWithFields(
				testutils.AMonthField(
					testutils.FieldID(fieldId),
					testutils.FieldName(fieldName),
				),
			),
		}, {
			name: "with week field",
			formDefinition: aFormWithFields(
				testutils.AWeekField(
					testutils.FieldID(fieldId),
					testutils.FieldName(fieldName),
				),
			),
		}, {
			name: "with quantity field",
			formDefinition: aFormWithFields(
				testutils.AQuantityField(
					testutils.FieldID(fieldId),
					testutils.FieldName(fieldName),
				),
			),
		}, {
			name: "with reference field",
			formDefinition: aFormWithFields(
				testutils.AReferenceField(aFormRef,
					testutils.FieldID(fieldId),
					testutils.FieldName(fieldName),
				),
			),
		}, {
			name: "with subform field",
			formDefinition: aFormWithFields(
				testutils.ASubFormField([]*types.FieldDefinition{
					testutils.ATextField(
						testutils.FieldID(fieldId),
						testutils.FieldName("sub field"),
					),
				}, testutils.FieldID(fieldId)),
			),
		}, {
			name: "with multiple sub forms",
			formDefinition: aFormWithFields(
				testutils.ASubFormField([]*types.FieldDefinition{
					testutils.ATextField(
						testutils.FieldID("subField1"),
						testutils.FieldName("sub field 1"),
					),
				}, testutils.FieldID("field1")),
				testutils.ASubFormField([]*types.FieldDefinition{
					testutils.ATextField(
						testutils.FieldID("subField2"),
						testutils.FieldName("sub field 2"),
					),
				}, testutils.FieldID("field2")),
			),
		}, {
			name: "with nested sub forms",
			formDefinition: aFormWithFields(
				testutils.ASubFormField([]*types.FieldDefinition{
					testutils.ASubFormField([]*types.FieldDefinition{
						testutils.ASubFormField([]*types.FieldDefinition{
							testutils.ATextField(
								testutils.FieldID("subField3"),
								testutils.FieldName("sub field 3"),
							),
						}, testutils.FieldID("field3")),
					}, testutils.FieldID("field2")),
				}, testutils.FieldID("field1")),
			),
		}, {
			name: "with single select field",
			formDefinition: aFormWithFields(
				testutils.ASingleSelectField([]*types.SelectOption{
					{
						ID:   "option 1",
						Name: "name 1",
					}, {
						ID:   "option 2",
						Name: "name 2",
					},
				},
					testutils.FieldName(fieldName),
					testutils.FieldID(fieldId)),
			),
		}, {
			name: "with formType",
			formDefinition: &types.FormDefinition{
				ID:         "formId",
				DatabaseID: "databaseId",
				Name:       "formName",
				Type:       types.RecipientFormType,
				Fields: []*types.FieldDefinition{
					{
						Name: "textField",
						ID:   "fieldId",
						FieldType: types.FieldType{
							Text: &types.FieldTypeText{},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flat, err := flattenForm(test.formDefinition)
			if !assert.NoError(t, err) {
				return
			}
			forms, err := flat.hydrateForms()
			if !assert.NoError(t, err) {
				return
			}
			if !assert.Len(t, forms, 1) {
				return
			}
			assert.Equal(t, test.formDefinition, forms[0])
		})
	}

}
