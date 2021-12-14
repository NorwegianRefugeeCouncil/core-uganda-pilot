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

func RecordValue(fieldID string, value types.StringOrArray) RecordOption {
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

type mockSubForm struct {
	*mockForm
	ownerId string
}

func (m mockSubForm) GetOwnerFormID() string {
	return m.ownerId
}

type mockForm struct {
	formId        string
	databaseId    string
	fields        types.FieldDefinitions
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

func ASubForm(ownerId string, options ...FormOption) *mockSubForm {
	r := &mockForm{}
	for _, option := range options {
		r = option(r)
	}
	return &mockSubForm{
		ownerId:  ownerId,
		mockForm: r,
	}
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
			aFormRef := types.FormRef{
				DatabaseID: uuid.NewV4().String(),
				FormID:     uuid.NewV4().String(),
			}
			return FieldTypeReference(aFormRef)(fieldDefinition)
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

func FieldKey(key bool) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.Key = true
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

func FieldTypeReference(formRef types.FormReference) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			Reference: &types.FieldTypeReference{
				DatabaseID: formRef.GetDatabaseID(),
				FormID:     formRef.GetFormID(),
			},
		}
		return fieldDefinition
	}
}

func FieldTypeSubForm(fields []*types.FieldDefinition) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			SubForm: &types.FieldTypeSubForm{
				Fields: fields,
			},
		}
		return fieldDefinition
	}
}

func FieldTypeSingleSelect(options []*types.SelectOption) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			SingleSelect: &types.FieldTypeSingleSelect{
				Options: options,
			},
		}
		return fieldDefinition
	}
}
func FieldTypeMultiSelect(options []*types.SelectOption) FieldOption {
	return func(fieldDefinition *types.FieldDefinition) *types.FieldDefinition {
		fieldDefinition.FieldType = types.FieldType{
			MultiSelect: &types.FieldTypeMultiSelect{
				Options: options,
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

func ATextField(options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeText()}
	opts = append(opts, options...)
	return AField(opts...)
}

func AMultilineTextField(options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeMultilineText()}
	opts = append(opts, options...)
	return AField(opts...)
}

func AMonthField(options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeMonth()}
	opts = append(opts, options...)
	return AField(opts...)
}

func ADateField(options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeDate()}
	opts = append(opts, options...)
	return AField(opts...)
}

func AWeekField(options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeWeek()}
	opts = append(opts, options...)
	return AField(opts...)
}

func AReferenceField(formRef types.FormReference, options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeReference(formRef)}
	opts = append(opts, options...)
	return AField(opts...)
}

func ASubFormField(fields []*types.FieldDefinition, options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeSubForm(fields)}
	opts = append(opts, options...)
	return AField(opts...)
}

func ASingleSelectField(selectOptions []*types.SelectOption, options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeSingleSelect(selectOptions)}
	opts = append(opts, options...)
	return AField(opts...)
}

func AMultiSelectField(selectOptions []*types.SelectOption, options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeMultiSelect(selectOptions)}
	opts = append(opts, options...)
	return AField(opts...)
}

func AQuantityField(options ...FieldOption) *types.FieldDefinition {
	opts := []FieldOption{FieldTypeQuantity()}
	opts = append(opts, options...)
	return AField(opts...)
}

func Fields(fields ...*types.FieldDefinition) []*types.FieldDefinition {
	return fields
}

func RecordForForm(form types.FormInterface) RecordOption {
	return func(r *types.Record) *types.Record {
		r.ID = uuid.NewV4().String()
		r.FormID = form.GetFormID()
		r.DatabaseID = form.GetDatabaseID()
		if r.Values == nil {
			r.Values = types.FieldValues{}
		}
		if _, ok := form.(types.SubFormInterface); ok {
			r.OwnerID = pointers.String(uuid.NewV4().String())
		}
		for _, field := range form.GetFields() {
			if field.FieldType.Date != nil {
				r.Values = append(r.Values, types.NewFieldStringValue(field.ID, "2020-01-01"))
			} else if field.FieldType.Month != nil {
				r.Values = append(r.Values, types.NewFieldStringValue(field.ID, "2020-01"))
			} else if field.FieldType.Week != nil {
				r.Values = append(r.Values, types.NewFieldStringValue(field.ID, "2020-W01"))
			} else if field.FieldType.Text != nil {
				r.Values = append(r.Values, types.NewFieldStringValue(field.ID, "abc"))
			} else if field.FieldType.MultilineText != nil {
				r.Values = append(r.Values, types.NewFieldStringValue(field.ID, "abc\ndef"))
			} else if field.FieldType.Reference != nil {
				r.Values = append(r.Values, types.NewFieldStringValue(field.ID, uuid.NewV4().String()))
			} else if field.FieldType.Quantity != nil {
				r.Values = append(r.Values, types.NewFieldStringValue(field.ID, "10"))
			} else if field.FieldType.SingleSelect != nil {
				r.Values = append(r.Values, types.NewFieldStringValue(field.ID, field.FieldType.SingleSelect.Options[0].ID))
			} else if field.FieldType.MultiSelect != nil {
				r.Values = append(r.Values, types.NewFieldArrayValue(field.ID, []string{field.FieldType.MultiSelect.Options[0].ID}))
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
