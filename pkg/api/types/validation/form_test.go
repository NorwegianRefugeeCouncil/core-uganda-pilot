package validation

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

var textFieldType = types.FieldType{
	Text: &types.FieldTypeText{},
}

var validTextField = &types.FieldDefinition{
	Name:      "My Field",
	FieldType: textFieldType,
}

const validFormName = "My Form"
const validFieldName = "My Field"

var validDatabaseID = uuid.NewV4().String()
var validFolderID = uuid.NewV4().String()

var validFields = types.FieldDefinitions{
	validTextField,
}

var formWithFields = func(fields types.FieldDefinitions) *types.FormDefinition {
	return &types.FormDefinition{
		Name:       validFormName,
		DatabaseID: validDatabaseID,
		Fields:     fields,
	}
}

var repeatFields = func(count int) types.FieldDefinitions {
	result := types.FieldDefinitions{}
	for i := 0; i < count; i++ {
		result = append(result, validTextField)
	}
	return result
}

var repeatOptions = func(count int) []*types.SelectOption {
	var result []*types.SelectOption
	for i := 0; i < count; i++ {
		result = append(result, &types.SelectOption{
			Name: strconv.Itoa(i),
		})
	}
	return result
}

func TestValidateForm(t *testing.T) {

	tests := []struct {
		name   string
		expect validation.ErrorList
		form   *types.FormDefinition
	}{
		{
			name:   "valid",
			expect: nil,
			form:   formWithFields(validFields),
		}, {
			name: "missing database id",
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("databaseId"), errDatabaseIdRequired),
			},
			form: &types.FormDefinition{
				Name:       validFormName,
				DatabaseID: "",
				Fields:     validFields,
			},
		}, {
			name: "form name with surrounding whitespaces",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("name"), " My Form ", errFormNameWhitespace),
			},
			form: &types.FormDefinition{
				Name:       " My Form ",
				DatabaseID: validDatabaseID,
				Fields:     validFields,
			},
		}, {
			name: "form name missing",
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("name"), errFormNameRequired),
			},
			form: &types.FormDefinition{
				Name:       "",
				DatabaseID: validDatabaseID,
				Fields:     validFields,
			},
		}, {
			name: "form name too long",
			expect: validation.ErrorList{
				validation.TooLong(validation.NewPath("name"), strings.Repeat("a", formNameMaxLength+1), formNameMaxLength),
			},
			form: &types.FormDefinition{
				Name:       strings.Repeat("a", formNameMaxLength+1),
				DatabaseID: validDatabaseID,
				Fields:     validFields,
			},
		}, {
			name: "form name too short",
			expect: validation.ErrorList{
				validation.TooShort(validation.NewPath("name"), strings.Repeat("a", formNameMinLength-1), formNameMinLength),
			},
			form: &types.FormDefinition{
				Name:       strings.Repeat("a", formNameMinLength-1),
				DatabaseID: validDatabaseID,
				Fields:     validFields,
			},
		}, {
			name: "bad database id",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("databaseId"), "abc", errInvalidUUID),
			},
			form: &types.FormDefinition{
				Name:       validFormName,
				DatabaseID: "abc",
				Fields:     validFields,
			},
		}, {
			name: "bad folder id",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("folderId"), "abc", errInvalidUUID),
			},
			form: &types.FormDefinition{
				Name:       validFormName,
				DatabaseID: validDatabaseID,
				FolderID:   "abc",
				Fields:     validFields,
			},
		}, {
			name:   "valid folder id",
			expect: nil,
			form: &types.FormDefinition{
				Name:       validFormName,
				DatabaseID: validDatabaseID,
				FolderID:   validFolderID,
				Fields:     validFields,
			},
		}, {
			name: "empty fields",
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("fields"), errFieldsRequired),
			},
			form: formWithFields(types.FieldDefinitions{}),
		}, {
			name: "too many fields",
			expect: validation.ErrorList{
				validation.TooMany(validation.NewPath("fields"), formMaxFieldCount+1, formMaxFieldCount),
			},
			form: formWithFields(repeatFields(formMaxFieldCount + 1)),
		}, {
			name: "field is key but not required",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].required"), false, errKeyFieldMustBeRequired),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:      validFieldName,
					Key:       true,
					Required:  false,
					FieldType: textFieldType,
				},
			}),
		}, {
			name: "sub form field cannot be required",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].required"), true, errSubFormCannotBeKeyOrRequiredField),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:     validFieldName,
					Key:      true,
					Required: true,
					FieldType: types.FieldType{
						SubForm: &types.FieldTypeSubForm{
							Fields: types.FieldDefinitions{
								validTextField,
							},
						},
					},
				},
			}),
		}, {
			name: "multiline text field cannot be key",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].key"), true, errMultiLineTextFieldCannotBeKeyField),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:     validFieldName,
					Key:      true,
					Required: true,
					FieldType: types.FieldType{
						MultilineText: &types.FieldTypeMultilineText{},
					},
				},
			}),
		}, {
			name: "field name cannot be empty",
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("fields[0].name"), errFieldNameRequired),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:      "",
					FieldType: textFieldType,
				},
			}),
		}, {
			name: "field name cannot be too long",
			expect: validation.ErrorList{
				validation.TooLong(validation.NewPath("fields[0].name"), strings.Repeat("a", fieldNameMaxLength+1), fieldNameMaxLength),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:      strings.Repeat("a", fieldNameMaxLength+1),
					FieldType: textFieldType,
				},
			}),
		}, {
			name: "field name cannot be too short",
			expect: validation.ErrorList{
				validation.TooShort(validation.NewPath("fields[0].name"), strings.Repeat("a", fieldNameMinLength-1), fieldNameMinLength),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:      strings.Repeat("a", fieldNameMinLength-1),
					FieldType: textFieldType,
				},
			}),
		}, {
			name: "field name cannot contain invalid characters",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].name"), "!!!", errFieldNameInvalid),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:      "!!!",
					FieldType: textFieldType,
				},
			}),
		}, {
			name: "field name cannot have surrounding whitespace",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].name"), " fieldName ", errFieldNameNoLeadingTrailingWhitespaces),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:      " fieldName ",
					FieldType: textFieldType,
				},
			}),
		}, {
			name: "field code cannot have invalid characters",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].code"), "!!!", errInvalidFieldCode),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:      validFieldName,
					Code:      "!!!",
					FieldType: textFieldType,
				},
			}),
		}, {
			name: "field code cannot be too long",
			expect: validation.ErrorList{
				validation.TooLong(validation.NewPath("fields[0].code"), strings.Repeat("a", fieldCodeMaxLength+1), fieldCodeMaxLength),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:      validFieldName,
					Code:      strings.Repeat("a", fieldCodeMaxLength+1),
					FieldType: textFieldType,
				},
			}),
		}, {
			name: "field must have field type",
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("fields[0].fieldType"), errOneFieldTypeRequired),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
				},
			}),
		}, {
			name: "field cannot have multiple field types",
			expect: validation.ErrorList{
				validation.TooLong(validation.NewPath("fields[0].fieldType"), fmt.Sprintf(errFieldTypesMultipleF, []types.FieldKind{types.FieldKindText, types.FieldKindMultilineText}), 1),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Text:          &types.FieldTypeText{},
						MultilineText: &types.FieldTypeMultilineText{},
					},
				},
			}),
		}, {
			name:   "field with valid text field",
			expect: nil,
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Text: &types.FieldTypeText{},
					},
				},
			}),
		}, {
			name:   "field with valid date field",
			expect: nil,
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Date: &types.FieldTypeDate{},
					},
				},
			}),
		}, {
			name:   "field with valid week field",
			expect: nil,
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Week: &types.FieldTypeWeek{},
					},
				},
			}),
		}, {
			name:   "field with valid month field",
			expect: nil,
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Month: &types.FieldTypeMonth{},
					},
				},
			}),
		}, {
			name:   "field with valid multiline text field",
			expect: nil,
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						MultilineText: &types.FieldTypeMultilineText{},
					},
				},
			}),
		}, {
			name:   "field with valid quantity field",
			expect: nil,
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Quantity: &types.FieldTypeQuantity{},
					},
				},
			}),
		}, {
			name:   "field with valid reference field",
			expect: nil,
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Reference: &types.FieldTypeReference{
							DatabaseID: uuid.NewV4().String(),
							FormID:     uuid.NewV4().String(),
						},
					},
				},
			}),
		}, {
			name: "reference field with invalid database id",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].fieldType.reference.databaseId"), "abc", errReferenceFieldDatabaseIdInvalid),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Reference: &types.FieldTypeReference{
							DatabaseID: "abc",
							FormID:     uuid.NewV4().String(),
						},
					},
				},
			}),
		}, {
			name: "reference field with empty database id",
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("fields[0].fieldType.reference.databaseId"), errReferenceFieldDatabaseIdRequired),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Reference: &types.FieldTypeReference{
							DatabaseID: "",
							FormID:     uuid.NewV4().String(),
						},
					},
				},
			}),
		}, {
			name: "reference field with invalid form id",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].fieldType.reference.formId"), "abc", errReferenceFieldFormIdInvalid),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Reference: &types.FieldTypeReference{
							DatabaseID: uuid.NewV4().String(),
							FormID:     "abc",
						},
					},
				},
			}),
		}, {
			name: "reference field with empty form id",
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("fields[0].fieldType.reference.formId"), errReferenceFieldFormIdRequired),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						Reference: &types.FieldTypeReference{
							DatabaseID: uuid.NewV4().String(),
							FormID:     "",
						},
					},
				},
			}),
		}, {
			name: "multi select field cannot be key",
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("fields[0].key"), true, errMultiSelectCannotBeKeyField),
			},
			form: formWithFields(types.FieldDefinitions{
				{
					Name:     validFieldName,
					Key:      true,
					Required: true,
					FieldType: types.FieldType{
						MultiSelect: &types.FieldTypeMultiSelect{
							Options: []*types.SelectOption{
								{Name: "option 1"},
								{Name: "option 2"},
							},
						},
					},
				},
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := ValidateForm(test.form)
			assert.Equal(t, test.expect, errs)
		})
	}

}

func TestValidateFormOptions(t *testing.T) {
	tests := []struct {
		name    string
		options []*types.SelectOption
		expect  func(fieldPath *validation.Path) validation.ErrorList
	}{
		{
			name: "valid",
			expect: func(fieldPath *validation.Path) validation.ErrorList {
				return nil
			},
			options: []*types.SelectOption{
				{
					Name: "option1",
					ID:   "option1",
				}, {
					Name: "option2",
					ID:   "option2",
				},
			},
		}, {
			name: "no options",
			expect: func(fieldPath *validation.Path) validation.ErrorList {
				return validation.ErrorList{
					validation.Required(
						fieldPath.Child("options"),
						errSelectOptionsRequired,
					),
				}
			},
			options: []*types.SelectOption{},
		}, {
			name: "duplicate option name",
			expect: func(fieldPath *validation.Path) validation.ErrorList {
				return validation.ErrorList{
					validation.Duplicate(
						fieldPath.
							Child("options").
							Index(1).
							Child("name"),
						"option 1",
					),
				}
			},
			options: []*types.SelectOption{
				{Name: "option 1"},
				{Name: "option 1"},
			},
		}, {
			name: "missing option name",
			expect: func(fieldPath *validation.Path) validation.ErrorList {
				return validation.ErrorList{
					validation.Required(
						fieldPath.
							Child("options").
							Index(0).
							Child("name"),
						errSelectOptionNameRequired,
					),
				}
			},
			options: []*types.SelectOption{
				{Name: ""},
			},
		}, {
			name: "invalid option name",
			expect: func(fieldPath *validation.Path) validation.ErrorList {
				return validation.ErrorList{
					validation.Invalid(
						fieldPath.Child("options").
							Index(0).
							Child("name"),
						"!!",
						errSelectOptionNameInvalid,
					),
				}
			},
			options: []*types.SelectOption{
				{Name: "!!"},
			},
		}, {
			name: "too many options",
			expect: func(fieldPath *validation.Path) validation.ErrorList {
				return validation.ErrorList{
					validation.TooMany(
						fieldPath.Child("options"),
						selectFieldMaxOptions+1,
						selectFieldMaxOptions,
					),
				}
			},
			options: repeatOptions(selectFieldMaxOptions + 1),
		},
	}
	t.Run("multiSelect", func(t *testing.T) {
		for _, test := range tests {
			form := formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						MultiSelect: &types.FieldTypeMultiSelect{
							Options: test.options,
						},
					},
				},
			})
			t.Run(test.name, func(t *testing.T) {
				errs := ValidateForm(form)
				path := validation.NewPath("fields").Index(0).Child("fieldType", "multiSelect")
				assert.Equal(t, test.expect(path), errs)
			})
		}
	})
	t.Run("singleSelect", func(t *testing.T) {
		for _, test := range tests {
			form := formWithFields(types.FieldDefinitions{
				{
					Name: validFieldName,
					FieldType: types.FieldType{
						SingleSelect: &types.FieldTypeSingleSelect{
							Options: test.options,
						},
					},
				},
			})
			t.Run(test.name, func(t *testing.T) {
				errs := ValidateForm(form)
				path := validation.NewPath("fields").Index(0).Child("fieldType", "singleSelect")
				assert.Equal(t, test.expect(path), errs)
			})
		}
	})
}

func TestValidateFieldNameRegex(t *testing.T) {
	valid := []string{
		"fieldName",
		"field name",
		"Field Name",
		"James Bond!- 007",
	}
	invalid := []string{
		" invalid ",
		" Field",
		"Field ",
		"Field  Field",
		"!Field",
		"    ",
	}
	for _, s := range valid {
		assert.True(t, fieldNameRegex.MatchString(s))
	}
	for _, s := range invalid {
		assert.False(t, fieldNameRegex.MatchString(s))
	}
}

func TestFieldCodeRegex(t *testing.T) {
	valid := []string{
		"code",
		"CODE",
		"code0",
	}
	invalid := []string{
		" ",
		" code",
		"code ",
		"code code",
		"!code",
		"0code",
	}
	for _, s := range valid {
		assert.True(t, fieldCodeRegex.MatchString(s))
	}
	for _, s := range invalid {
		assert.False(t, fieldCodeRegex.MatchString(s))
	}
}
