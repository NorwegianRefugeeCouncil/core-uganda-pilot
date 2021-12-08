package types

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFindSubForm3(t *testing.T) {

	tests := []struct {
		name        string
		id          string
		form        FormInterface
		expectErr   bool
		expectFound bool
	}{
		{
			name:        "FormDefinition empty",
			expectFound: false,
			form: &FormDefinition{
				ID: "myform",
			},
		}, {
			name:        "FormDefinition with subForm",
			id:          "subForm",
			expectFound: true,
			form: &FormDefinition{
				Fields: FieldDefinitions{
					{
						ID: "subForm",
						FieldType: FieldType{
							SubForm: &FieldTypeSubForm{},
						},
					},
				},
			},
		}, {
			name:        "FormDefinition with nested sub form",
			id:          "subForm",
			expectFound: true,
			form: &FormDefinition{
				Fields: FieldDefinitions{
					{
						ID: "subForm",
						FieldType: FieldType{
							SubForm: &FieldTypeSubForm{
								Fields: FieldDefinitions{
									{
										ID: "nestedSubForm",
										FieldType: FieldType{
											SubForm: &FieldTypeSubForm{
												Fields: FieldDefinitions{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, {
			name:        "FormDefinition with wrong field type",
			id:          "someField",
			expectFound: false,
			expectErr:   true,
			form: &FormDefinition{
				Fields: FieldDefinitions{
					{
						ID: "someField",
						FieldType: FieldType{
							Text: &FieldTypeText{},
						},
					},
				},
			},
		}, {
			name:        "subFormInterface child",
			expectFound: true,
			id:          "subForm",
			form: &subFormInterface{
				id: "subFormInterface",
				fields: FieldDefinitions{
					{
						ID: "subForm",
						FieldType: FieldType{
							SubForm: &FieldTypeSubForm{
								Fields: FieldDefinitions{},
							},
						},
					},
				},
				databaseId:  "databaseId",
				ownerFormID: "",
			},
		}, {
			name:        "subFormInterface nested child",
			expectFound: true,
			id:          "nestedSubForm",
			form: &subFormInterface{
				id: "subFormInterface",
				fields: FieldDefinitions{
					{
						ID: "subForm",
						FieldType: FieldType{
							SubForm: &FieldTypeSubForm{
								Fields: FieldDefinitions{
									{
										ID: "nestedSubForm",
										FieldType: FieldType{
											SubForm: &FieldTypeSubForm{
												Fields: FieldDefinitions{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, {
			name:      "subFormInterface missing field type",
			expectErr: true,
			id:        "missingType",
			form: &subFormInterface{
				id: "subFormInterface",
				fields: FieldDefinitions{
					{
						ID: "missingType",
					},
				},
			},
		}, {
			name:      "subFormInterface missing nested field type",
			expectErr: true,
			id:        "missing-field-type",
			form: &subFormInterface{
				id: "subFormInterface",
				fields: FieldDefinitions{
					{
						ID: "subForm",
						FieldType: FieldType{
							SubForm: &FieldTypeSubForm{
								Fields: FieldDefinitions{
									{
										ID: "missing-field-type",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.form.FindSubForm(test.id)
			if test.expectErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			val := reflect.ValueOf(got)
			found := val.Kind() == reflect.Ptr && !val.IsNil()
			assert.Equal(t, test.expectFound, found)
			if found && test.expectFound {
				assert.Equal(t, test.id, got.GetFormID())
			}
		})
	}

}
