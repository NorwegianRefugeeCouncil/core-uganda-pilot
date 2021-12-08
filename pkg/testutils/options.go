package testutils

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/satori/go.uuid"
)

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

func RecordID(recordId string) RecordOption {
	return func(record *types.Record) *types.Record {
		record.ID = recordId
		return record
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

func RecordValues(values types.FieldValues) RecordOption {
	return func(record *types.Record) *types.Record {

		if values == nil {
			record.Values = nil
			return record
		}

		var v = types.FieldValues{}
		for _, value := range values {
			v = append(v, value)
		}
		record.Values = v
		return record
	}
}

func RecordValue(fieldID string, value *string) RecordOption {
	return func(record *types.Record) *types.Record {
		val := types.FieldValue{
			FieldID: fieldID,
			Value:   value,
		}
		for i, fieldValue := range record.Values {
			if fieldValue.FieldID == fieldID {
				record.Values[i] = val
				return record
			}
		}
		record.Values = append(record.Values, val)
		return record
	}
}
func RecordOmitValue(fieldID string) RecordOption {
	return func(record *types.Record) *types.Record {
		res := types.FieldValues{}
		for _, value := range record.Values {
			if value.FieldID != fieldID {
				res = append(res, value)
			}
		}
		record.Values = res
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

func (m mockForm) IsSubForm() bool {
	return m.hasOwner
}

func (m mockForm) FindSubForm(subFormId string) (types.FormInterface, error) {
	return m.findSubFormFn()
}

var _ types.FormInterface = &mockForm{}

type FormOption func(form *mockForm) *mockForm

func FormOptions(options ...FormOption) FormOption {
	return func(form *mockForm) *mockForm {
		walk := form
		for _, option := range options {
			walk = option(walk)
		}
		return walk
	}
}

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

func FieldTypeKind(kind types.FieldKind) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		switch kind {
		case types.FieldKindUnknown:
			fieldDefinition.FieldType = types.FieldType{}
		case types.FieldKindText:
			return FieldTypeText()(fieldDefinition)
		case types.FieldKindSubForm:
		case types.FieldKindReference:
			return FieldTypeReference()(fieldDefinition)
		case types.FieldKindMultilineText:
			return FieldTypeMultilineText()(fieldDefinition)
		case types.FieldKindDate:
			return FieldTypeDate()(fieldDefinition)
		case types.FieldKindQuantity:
			return FieldTypeQuantity()(fieldDefinition)
		case types.FieldKindMonth:
			return FieldTypeMonth()(fieldDefinition)
		case types.FieldKindWeek:
			return FieldTypeWeek()(fieldDefinition)
		case types.FieldKindSingleSelect:
			// todo
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

func FieldTypeReference() FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			Reference: &types.FieldTypeReference{
				DatabaseID: uuid.NewV4().String(),
				FormID:     uuid.NewV4().String(),
			},
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

func FieldTypeWeek() FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			Week: &types.FieldTypeWeek{},
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
			r.Values = types.FieldValues{}
		}
		if form.IsSubForm() {
			r.OwnerID = pointers.String(uuid.NewV4().String())
		}
		for _, field := range form.GetFields() {
			if field.FieldType.Date != nil {
				r.Values = append(r.Values, types.FieldValue{
					FieldID: field.ID,
					Value:   pointers.String("2020-01-01"),
				})
			} else if field.FieldType.Month != nil {
				r.Values = append(r.Values, types.FieldValue{
					FieldID: field.ID,
					Value:   pointers.String("2020-01"),
				})
			} else if field.FieldType.Week != nil {
				r.Values = append(r.Values, types.FieldValue{
					FieldID: field.ID,
					Value:   pointers.String("2020-W01"),
				})
			} else if field.FieldType.Text != nil {
				r.Values = append(r.Values, types.FieldValue{
					FieldID: field.ID,
					Value:   pointers.String("abc"),
				})
			} else if field.FieldType.MultilineText != nil {
				r.Values = append(r.Values, types.FieldValue{
					FieldID: field.ID,
					Value:   pointers.String("abc\ndef"),
				})
			} else if field.FieldType.Reference != nil {
				r.Values = append(r.Values, types.FieldValue{
					FieldID: field.ID,
					Value:   pointers.String(uuid.NewV4().String()),
				})
			} else if field.FieldType.Quantity != nil {
				r.Values = append(r.Values, types.FieldValue{
					FieldID: field.ID,
					Value:   pointers.String("10"),
				})
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
