package validation

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateRecord(t *testing.T) {

	var (
		formId     = uuid.NewV4().String()
		fieldId    = uuid.NewV4().String()
		databaseId = uuid.NewV4().String()
	)

	aTextForm := func(options ...FormOption) types.FormInterface {
		opts := []FormOption{
			FormID(formId),
			FormDatabaseID(databaseId),
			FormField(&types.FieldDefinition{
				ID: fieldId,
				FieldType: types.FieldType{
					Text: &types.FieldTypeText{},
				},
			}),
		}
		opts = append(opts, options...)
		f := AForm(opts...)
		return f
	}

	textForm := aTextForm()

	tests := []struct {
		name          string
		recordOptions RecordOption
		form          types.FormInterface
		expect        validation.ErrorList
	}{
		{
			name:   "valid",
			form:   textForm,
			expect: nil,
		}, {
			name:          "missing form id",
			form:          aTextForm(),
			recordOptions: RecordFormID(""),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("formId"), errRecordFormIdRequired),
			},
		}, {
			name:          "invalid form id",
			form:          textForm,
			recordOptions: RecordFormID("bla"),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("formId"), "bla", errRecordInvalidFormId),
			},
		}, {
			name:          "missing database id",
			form:          textForm,
			recordOptions: RecordDatabaseID(""),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("databaseId"), errRecordDatabaseIdRequired),
			},
		}, {
			name:          "invalid database id",
			form:          textForm,
			recordOptions: RecordDatabaseID("bla"),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("databaseId"), "bla", errRecordInvalidDatabaseId),
			},
		}, {
			name:          "missing ownerId",
			form:          aTextForm(FormHasOwner(true)),
			recordOptions: RecordOwnerID(nil),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("ownerId"), errRecordOwnerIdRequired),
			},
		}, {
			name:          "empty ownerId",
			form:          aTextForm(FormHasOwner(true)),
			recordOptions: RecordOwnerID(pointers.String("")),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("ownerId"), errRecordOwnerIdRequired),
			},
		}, {
			name:          "invalid ownerId",
			form:          aTextForm(FormHasOwner(true)),
			recordOptions: RecordOwnerID(pointers.String("abc")),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("ownerId"), "abc", errRecordInvalidOwnerID),
			},
		}, {
			name:          "nil values",
			form:          aTextForm(),
			recordOptions: RecordValues(nil),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("values"), errRecordValuesRequired),
			},
		}, {
			name:          "missing field type",
			form:          aTextForm(FormField(AField(FieldID("someField")))),
			recordOptions: RecordValue("bla", "snip"),
			expect: validation.ErrorList{
				validation.InternalError(validation.NewPath("values"), errors.New("failed to get field kind")),
			},
		}, {
			name:          "extraneous field",
			form:          aTextForm(),
			recordOptions: RecordValue("bla", "snip"),
			expect: validation.ErrorList{
				validation.NotSupported(validation.NewPath("values").Key("bla"), "bla", []string{fieldId}),
			},
		}, {
			name: "missing required field",
			form: aTextForm(
				FormField(AField(FieldID("requiredField"), FieldTypeText(), FieldRequired(true))),
			),
			recordOptions: RecordOmitValue("requiredField"),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("values").Key("requiredField"), errFieldValueRequired),
			},
		}, {
			name: "zero-valued required field",
			form: aTextForm(
				FormField(AField(FieldID("requiredField"), FieldTypeText(), FieldRequired(true))),
			),
			recordOptions: RecordValue("requiredField", ""),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("values").Key("requiredField"), errFieldValueRequired),
			},
		}, {
			name: "missing optional field",
			form: aTextForm(
				FormField(AField(FieldID("required"), FieldTypeText(), FieldRequired(false))),
			),
			expect: nil,
		}, {
			name: "invalid text field",
			form: aTextForm(
				FormField(AField(FieldID("textField"), FieldTypeText())),
			),
			recordOptions: RecordValue("textField", 123),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("values").Key("textField"), 123, fmt.Sprintf(errInvalidFieldValueTypeF, "", 123)),
			},
		}, {
			name: "invalid multiline field",
			form: aTextForm(
				FormField(AField(FieldID("multiLineTextField"), FieldTypeMultilineText())),
			),
			recordOptions: RecordValue("multiLineTextField", 123),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("values").Key("multiLineTextField"), 123, fmt.Sprintf(errInvalidFieldValueTypeF, "", 123)),
			},
		}, {
			name: "date field",
			form: aTextForm(
				FormField(AField(FieldID("dateField"), FieldTypeDate())),
			),
			expect: nil,
		}, {
			name: "invalid date field",
			form: aTextForm(
				FormField(AField(FieldID("dateField"), FieldTypeDate())),
			),
			recordOptions: RecordValue("dateField", "someValue"),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("values").Key("dateField"), "someValue", errRecordInvalidDate),
			},
		}, {
			name: "date field wrong type",
			form: aTextForm(
				FormField(AField(FieldID("dateField"), FieldTypeDate())),
			),
			recordOptions: RecordValue("dateField", 123),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("values").Key("dateField"), 123, fmt.Sprintf(errInvalidFieldValueTypeF, "", 123)),
			},
		}, {
			name: "required empty date field",
			form: aTextForm(
				FormField(AField(FieldID("dateField"), FieldTypeDate(), FieldRequired(true))),
			),
			recordOptions: RecordValue("dateField", ""),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("values").Key("dateField"), errFieldValueRequired),
			},
		}, {
			name: "month field",
			form: aTextForm(
				FormField(AField(FieldID("dateField"), FieldTypeDate())),
			),
			expect: nil,
		}, {
			name: "invalid month field",
			form: aTextForm(
				FormField(AField(FieldID("monthField"), FieldTypeMonth())),
			),
			recordOptions: RecordValue("monthField", "someValue"),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("values").Key("monthField"), "someValue", errRecordInvalidMonth),
			},
		}, {
			name: "month field wrong type",
			form: aTextForm(
				FormField(AField(FieldID("monthField"), FieldTypeMonth())),
			),
			recordOptions: RecordValue("monthField", 123),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("values").Key("monthField"), 123, fmt.Sprintf(errInvalidFieldValueTypeF, "", 123)),
			},
		}, {
			name: "required empty month field",
			form: aTextForm(
				FormField(AField(FieldID("monthField"), FieldTypeMonth(), FieldRequired(true))),
			),
			recordOptions: RecordValue("monthField", ""),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("values").Key("monthField"), errFieldValueRequired),
			},
		}, {
			name: "int field",
			form: aTextForm(
				FormField(AField(FieldID("dateField"), FieldTypeDate())),
			),
			expect: nil,
		}, {
			name: "invalid int field",
			form: aTextForm(
				FormField(AField(FieldID("quantityField"), FieldTypeQuantity())),
			),
			recordOptions: RecordValue("quantityField", "someValue"),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("values").Key("quantityField"), "someValue", errRecordInvalidQuantity),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts := []RecordOption{RecordForForm(test.form)}
			if test.recordOptions != nil {
				opts = append(opts, test.recordOptions)
			}
			rec := ARecord(opts...)
			assert.Equal(t, test.expect, ValidateRecord(rec, test.form))
		})
	}

}

type RecordOption func(record *types.Record) *types.Record

func RecordOptions(options ...RecordOption) RecordOption {
	return func(record *types.Record) *types.Record {
		r := record
		for _, option := range options {
			r = option(r)
		}
		return r
	}
}

func RecordFormID(formID string) RecordOption {
	return func(record *types.Record) *types.Record {
		record.FormID = formID
		return record
	}
}

func RecordDatabaseID(databaseID string) RecordOption {
	return func(record *types.Record) *types.Record {
		record.DatabaseID = databaseID
		return record
	}
}

func RecordOwnerID(ownerID *string) RecordOption {
	return func(record *types.Record) *types.Record {
		record.OwnerID = ownerID
		return record
	}
}

func RecordValues(values map[string]interface{}) RecordOption {
	return func(record *types.Record) *types.Record {
		if values == nil {
			record.Values = nil
			return record
		}
		var v = make(map[string]interface{})
		for key, value := range values {
			v[key] = value
		}
		record.Values = v
		return record
	}
}

func RecordValue(key string, value interface{}) RecordOption {
	return func(record *types.Record) *types.Record {
		if record.Values == nil {
			record.Values = map[string]interface{}{}
		}
		record.Values[key] = value
		return record
	}
}
func RecordOmitValue(key string) RecordOption {
	return func(record *types.Record) *types.Record {
		if record.Values != nil {
			delete(record.Values, key)
		}
		return record
	}
}

func ARecord(options ...RecordOption) *types.Record {
	r := &types.Record{}
	for _, option := range options {
		r = option(r)
	}
	return r
}

type mockForm struct {
	formId        string
	databaseId    string
	fields        types.FieldDefinitions
	owner         types.FormInterface
	hasOwner      bool
	findSubFormFn func() (types.FormInterface, error)
}

func (m mockForm) GetFormID() string {
	return m.formId
}

func (m mockForm) GetDatabaseID() string {
	return m.databaseId
}

func (m mockForm) GetFields() types.FieldDefinitions {
	return m.fields
}

func (m mockForm) GetOwner() types.FormInterface {
	return m.owner
}

func (m mockForm) HasOwner() bool {
	return m.hasOwner
}

func (m mockForm) FindSubForm(subFormId string) (types.FormInterface, error) {
	return m.findSubFormFn()
}

var _ types.FormInterface = &mockForm{}

type FormOption func(form *mockForm) *mockForm

func FormID(id string) FormOption {
	return func(form *mockForm) *mockForm {
		form.formId = id
		return form
	}
}

func FormDatabaseID(databaseId string) FormOption {
	return func(form *mockForm) *mockForm {
		form.databaseId = databaseId
		return form
	}
}

func FormHasOwner(hasOwner bool) FormOption {
	return func(form *mockForm) *mockForm {
		form.hasOwner = hasOwner
		return form
	}
}

func FormOwner(owner types.FormInterface) FormOption {
	return func(form *mockForm) *mockForm {
		form.owner = owner
		return form
	}
}

func FormField(field *types.FieldDefinition) FormOption {
	return func(form *mockForm) *mockForm {
		form.fields = append(form.fields, field)
		return form
	}
}

func FormFields(fields ...*types.FieldDefinition) FormOption {
	return func(form *mockForm) *mockForm {
		form.fields = fields
		return form
	}
}

func AForm(options ...FormOption) *mockForm {
	r := &mockForm{}
	for _, option := range options {
		r = option(r)
	}
	return r
}

type FieldOption func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition

func FieldID(fieldID string) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.ID = fieldID
		return fieldDefinition
	}
}

func FieldName(fieldName string) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.Name = fieldName
		return fieldDefinition
	}
}

func FieldTypeText() FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			Text: &types.FieldTypeText{},
		}
		return fieldDefinition
	}
}

func FieldTypeQuantity() FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			Quantity: &types.FieldTypeQuantity{},
		}
		return fieldDefinition
	}
}

func FieldRequired(required bool) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.Required = required
		return fieldDefinition
	}
}

func FieldTypeMultilineText() FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			MultilineText: &types.FieldTypeMultilineText{},
		}
		return fieldDefinition
	}
}

func FieldTypeDate() FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			Date: &types.FieldTypeDate{},
		}
		return fieldDefinition
	}
}

func FieldTypeMonth() FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			Month: &types.FieldTypeMonth{},
		}
		return fieldDefinition
	}
}

func AField(options ...FieldOption) *types.FieldDefinition {
	f := &types.FieldDefinition{}
	for _, option := range options {
		f = option(f)
	}
	return f
}

func RecordForForm(form types.FormInterface) RecordOption {
	return func(r *types.Record) *types.Record {
		r.ID = uuid.NewV4().String()
		r.FormID = form.GetFormID()
		r.DatabaseID = form.GetDatabaseID()
		if r.Values == nil {
			r.Values = map[string]interface{}{}
		}
		if form.HasOwner() {
			r.OwnerID = pointers.String(uuid.NewV4().String())
		}
		for _, field := range form.GetFields() {
			if field.FieldType.Date != nil {
				r.Values[field.ID] = "2020-01-01"
			} else if field.FieldType.Month != nil {
				r.Values[field.ID] = "2020-01"
			} else if field.FieldType.Text != nil {
				r.Values[field.ID] = "abc"
			} else if field.FieldType.MultilineText != nil {
				r.Values[field.ID] = "abc\ndef"
			} else if field.FieldType.Reference != nil {
				r.Values[field.ID] = uuid.NewV4().String()
			} else if field.FieldType.Quantity != nil {
				r.Values[field.ID] = 3
			}
		}
		return r
	}
}

func ARecordForForm(form types.FormInterface, options ...RecordOption) *types.Record {
	opts := []RecordOption{
		RecordForForm(form),
	}
	opts = append(opts, options...)
	return ARecord(opts...)
}
