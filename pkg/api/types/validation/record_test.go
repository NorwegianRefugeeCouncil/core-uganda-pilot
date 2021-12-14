package validation

import (
	"errors"
	"github.com/nrc-no/core/pkg/api/types"
	tu "github.com/nrc-no/core/pkg/testutils"
	"github.com/nrc-no/core/pkg/utils/pointers"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateRecord(t *testing.T) {

	var (
		formId      = uuid.NewV4().String()
		fieldId     = uuid.NewV4().String()
		databaseId  = uuid.NewV4().String()
		ownerFormId = uuid.NewV4().String()
	)

	textFormOpts := []tu.FormOption{
		tu.FormID(formId),
		tu.FormDatabaseID(databaseId),
		tu.FormField(&types.FieldDefinition{
			ID: fieldId,
			FieldType: types.FieldType{
				Text: &types.FieldTypeText{},
			},
		}),
	}

	aTextForm := func(options ...tu.FormOption) types.FormInterface {
		f := tu.AForm(append(textFormOpts, options...)...)
		return f
	}

	aTextSubForm := func(options ...tu.FormOption) types.FormInterface {
		f := tu.ASubForm(ownerFormId, append(textFormOpts, options...)...)
		return f
	}

	textForm := aTextForm()

	valuePath := validation.NewPath("values")
	firstFieldPath := valuePath.Index(1)
	firstFieldValuePath := firstFieldPath.Child("value")

	aFormRef := types.FormRef{
		DatabaseID: uuid.NewV4().String(),
		FormID:     uuid.NewV4().String(),
	}

	selectOptions := []*types.SelectOption{
		{
			ID:   "option1",
			Name: "Option 1",
		}, {
			ID:   "option2",
			Name: "Option 2",
		},
	}

	tests := []struct {
		name          string
		recordOptions tu.RecordOption
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
			recordOptions: tu.RecordFormID(""),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("formId"), errRecordFormIdRequired),
			},
		}, {
			name:          "invalid form id",
			form:          textForm,
			recordOptions: tu.RecordFormID("bla"),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("formId"), "bla", errRecordInvalidFormId),
			},
		}, {
			name:          "missing database id",
			form:          textForm,
			recordOptions: tu.RecordDatabaseID(""),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("databaseId"), errRecordDatabaseIdRequired),
			},
		}, {
			name:          "invalid database id",
			form:          textForm,
			recordOptions: tu.RecordDatabaseID("bla"),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("databaseId"), "bla", errRecordInvalidDatabaseId),
			},
		}, {
			name:          "missing ownerId",
			form:          aTextSubForm(),
			recordOptions: tu.RecordOwnerID(nil),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("ownerId"), errRecordOwnerIdRequired),
			},
		}, {
			name:          "empty ownerId",
			form:          aTextSubForm(),
			recordOptions: tu.RecordOwnerID(pointers.String("")),
			expect: validation.ErrorList{
				validation.Required(validation.NewPath("ownerId"), errRecordOwnerIdRequired),
			},
		}, {
			name:          "invalid ownerId",
			form:          aTextSubForm(),
			recordOptions: tu.RecordOwnerID(pointers.String("abc")),
			expect: validation.ErrorList{
				validation.Invalid(validation.NewPath("ownerId"), "abc", errRecordInvalidOwnerID),
			},
		}, {
			name:          "nil values",
			form:          aTextForm(),
			recordOptions: tu.RecordValues(nil),
			expect: validation.ErrorList{
				validation.Required(valuePath, errRecordValuesRequired),
			},
		}, {
			name:          "missing field type",
			form:          aTextForm(tu.FormField(tu.AField(tu.FieldID("someField")))),
			recordOptions: tu.RecordValue("bla", types.NewStringValue("snip")),
			expect: validation.ErrorList{
				validation.InternalError(valuePath, errors.New("failed to get field kind")),
			},
		}, {
			name:          "extraneous field",
			form:          aTextForm(),
			recordOptions: tu.RecordValue("bla", types.NewStringValue("snip")),
			expect: validation.ErrorList{
				validation.NotSupported(firstFieldPath, "bla", []string{fieldId}),
			},
		}, {
			name: "missing required field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("requiredField"), tu.FieldTypeText(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordOmitValue("requiredField"),
			expect: validation.ErrorList{
				validation.Required(valuePath, errFieldValueRequired),
			},
		}, {
			name: "zero-valued required field",
			form: aTextForm(
				tu.FormField(
					tu.AField(
						tu.FieldID("requiredField"),
						tu.FieldTypeText(),
						tu.FieldRequired(true),
					),
				),
			),
			recordOptions: tu.RecordValue("requiredField", types.NewStringValue("")),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "missing optional field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("required"), tu.FieldTypeText(), tu.FieldRequired(false))),
			),
			expect: nil,
		}, {
			name: "required text field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("textField"), tu.FieldTypeText(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("textField", types.NewNullValue()),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "optional text field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("textField"), tu.FieldTypeText())),
			),
			recordOptions: tu.RecordValue("textField", types.NewNullValue()),
			expect:        nil,
		}, {
			name: "required multiline text field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("textField"), tu.FieldTypeMultilineText(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("textField", types.NewNullValue()),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "optional multiline text field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("textField"), tu.FieldTypeMultilineText())),
			),
			recordOptions: tu.RecordValue("textField", types.NewNullValue()),
			expect:        nil,
		}, {
			name: "date field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("dateField"), tu.FieldTypeDate())),
			),
			expect: nil,
		}, {
			name: "invalid date field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("dateField"), tu.FieldTypeDate())),
			),
			recordOptions: tu.RecordValue("dateField", types.NewStringValue("someValue")),
			expect: validation.ErrorList{
				validation.Invalid(firstFieldValuePath, types.NewStringValue("someValue"), errRecordInvalidDate),
			},
		}, {
			name: "required empty date field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("dateField"), tu.FieldTypeDate(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("dateField", types.NewStringValue("")),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "required date field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("dateField"), tu.FieldTypeDate(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("dateField", types.NewNullValue()),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "optional date field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("dateField"), tu.FieldTypeDate())),
			),
			recordOptions: tu.RecordValue("dateField", types.NewNullValue()),
			expect:        nil,
		}, {
			name: "month field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("monthField"), tu.FieldTypeMonth())),
			),
			expect: nil,
		}, {
			name: "invalid month field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("monthField"), tu.FieldTypeMonth())),
			),
			recordOptions: tu.RecordValue("monthField", types.NewStringValue("someValue")),
			expect: validation.ErrorList{
				validation.Invalid(firstFieldValuePath, types.NewStringValue("someValue"), errRecordInvalidMonth),
			},
		}, {
			name: "required empty month field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("monthField"), tu.FieldTypeMonth(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("monthField", types.NewStringValue("")),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "required month field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("monthField"), tu.FieldTypeMonth(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("monthField", types.NewNullValue()),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "optional month field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("monthField"), tu.FieldTypeMonth())),
			),
			recordOptions: tu.RecordValue("monthField", types.NewNullValue()),
			expect:        nil,
		}, {
			name: "week field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("weekField"), tu.FieldTypeWeek())),
			),
			expect: nil,
		}, {
			name: "invalid week field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("weekField"), tu.FieldTypeWeek())),
			),
			recordOptions: tu.RecordValue("weekField", types.NewStringValue("someValue")),
			expect: validation.ErrorList{
				validation.Invalid(firstFieldValuePath, types.NewStringValue("someValue"), errRecordInvalidWeek),
			},
		}, {
			name: "required empty week field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("weekField"), tu.FieldTypeWeek(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("weekField", types.NewStringValue("")),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "required week field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("weekField"), tu.FieldTypeWeek(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("weekField", types.NewNullValue()),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "optional week field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("weekField"), tu.FieldTypeWeek())),
			),
			recordOptions: tu.RecordValue("weekField", types.NewNullValue()),
			expect:        nil,
		}, {
			name: "quantity field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("dateField"), tu.FieldTypeDate())),
			),
			expect: nil,
		}, {
			name: "invalid quantity field value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("quantityField"), tu.FieldTypeQuantity())),
			),
			recordOptions: tu.RecordValue("quantityField", types.NewStringValue("someValue")),
			expect: validation.ErrorList{
				validation.Invalid(firstFieldValuePath, types.NewStringValue("someValue"), errRecordInvalidQuantity),
			},
		}, {
			name: "required quantity field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("quantityField"), tu.FieldTypeQuantity(), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("quantityField", types.NewNullValue()),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "optional quantity field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("quantityField"), tu.FieldTypeQuantity())),
			),
			recordOptions: tu.RecordValue("quantityField", types.NewNullValue()),
			expect:        nil,
		}, {
			name: "reference field",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("dateField"), tu.FieldTypeReference(aFormRef))),
			),
			expect: nil,
		}, {
			name: "invalid reference field value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("referenceField"), tu.FieldTypeReference(aFormRef))),
			),
			recordOptions: tu.RecordValue("referenceField", types.NewStringValue("someValue")),
			expect: validation.ErrorList{
				validation.Invalid(firstFieldValuePath, types.NewStringValue("someValue"), errRecordInvalidReferenceUid),
			},
		}, {
			name: "required reference field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("referenceField"), tu.FieldTypeReference(aFormRef), tu.FieldRequired(true))),
			),
			recordOptions: tu.RecordValue("referenceField", types.NewNullValue()),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "optional reference field with nil value",
			form: aTextForm(
				tu.FormField(tu.AField(tu.FieldID("referenceField"), tu.FieldTypeReference(aFormRef))),
			),
			recordOptions: tu.RecordValue("referenceField", types.NewNullValue()),
			expect:        nil,
		}, {
			name: "single select field",
			form: aTextForm(
				tu.FormField(
					tu.ASingleSelectField(
						selectOptions,
						tu.FieldID("singleSelectField"),
					),
				),
			),
			expect: nil,
		}, {
			name: "single select field with unknown value",
			form: aTextForm(
				tu.FormField(
					tu.ASingleSelectField(
						selectOptions,
						tu.FieldID("singleSelectField"),
					),
				),
			),
			recordOptions: tu.RecordValue("singleSelectField", types.NewStringValue("someRandomValue")),
			expect: validation.ErrorList{
				validation.NotSupported(firstFieldValuePath, "someRandomValue", []string{
					"option1",
					"option2",
				}),
			},
		}, {
			name: "required single select field with nil value",
			form: aTextForm(
				tu.FormField(
					tu.ASingleSelectField(
						selectOptions,
						tu.FieldID("singleSelectField"),
						tu.FieldRequired(true),
					),
				),
			),
			recordOptions: tu.RecordValue("singleSelectField", types.NewNullValue()),
			expect: validation.ErrorList{
				validation.Required(firstFieldValuePath, errFieldValueRequired),
			},
		}, {
			name: "optional single select field with nil value",
			form: aTextForm(
				tu.FormField(
					tu.ASingleSelectField(
						selectOptions,
						tu.FieldID("singleSelectField"),
					),
				),
			),
			recordOptions: tu.RecordValue("singleSelectField", types.NewNullValue()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts := []tu.RecordOption{tu.RecordForForm(test.form)}
			if test.recordOptions != nil {
				opts = append(opts, test.recordOptions)
			}
			rec := tu.ARecord(opts...)
			got := ValidateRecord(rec, test.form)
			assert.Equal(t, test.expect, got)
		})
	}

}
