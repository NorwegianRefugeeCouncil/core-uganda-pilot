package validation

import (
	"strconv"
	"strings"
	"time"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/dates"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
)

const (
	errRecordInvalidDatabaseId   = "Invalid database ID"
	errRecordDatabaseIdRequired  = "Database ID is required"
	errRecordInvalidFormId       = "Invalid form ID"
	errRecordFormIdRequired      = "Record form ID is required"
	errRecordOwnerIdRequired     = "Record owner ID is required"
	errRecordInvalidOwnerID      = "Record owner ID is invalid"
	errRecordValuesRequired      = "Record values are required"
	errRecordInvalidDate         = "Invalid date. Expected YYYY-mm-DD"
	errRecordInvalidMonth        = "Invalid date. Expected YYYY-mm"
	errRecordInvalidWeek         = "Invalid date. Expected YYYY-Www"
	errRecordInvalidQuantity     = "Invalid quantity"
	errRecordInvalidReferenceUid = "Invalid reference"
	errFieldValueRequired        = "Field value is required"
	errFieldValueMustBeString    = "Field value must be a string"
	errFieldValueMustBeArray     = "Field value must be a string array"
)

// supportedRecordFieldKinds are the types of field for which a Record can specify values for
// in the types.Record values. For example, a types.FieldKindSubForm should not have a
// value in the types.Record values.
var supportedRecordFieldKinds = []types.FieldKind{
	types.FieldKindText,
	types.FieldKindReference,
	types.FieldKindMultilineText,
	types.FieldKindDate,
	types.FieldKindQuantity,
	types.FieldKindMonth,
	types.FieldKindWeek,
	types.FieldKindSingleSelect,
	types.FieldKindMultiSelect,
	types.FieldKindBoolean,
}

// supportedRecordFieldKindMap is a map of the supportedRecordFieldKinds for faster lookup
var supportedRecordFieldKindMap = map[types.FieldKind]struct{}{}

// This is a list of the names of supportedRecordFieldKinds for friendly error messages
var supportedFieldKindNames []string

func init() {
	// initialises the supportedFieldKindNames list
	for _, kind := range supportedRecordFieldKinds {
		supportedFieldKindNames = append(supportedFieldKindNames, kind.String())
		supportedRecordFieldKindMap[kind] = struct{}{}
	}
}

func ValidateRecord(record *types.Record, form types.FormInterface) validation.ErrorList {
	var result validation.ErrorList

	databaseIdPath := validation.NewPath("databaseId")
	formIdPath := validation.NewPath("formId")
	ownerIdPath := validation.NewPath("ownerId")
	valuesPath := validation.NewPath("values")

	if len(record.DatabaseID) == 0 {
		result = append(result, validation.Required(databaseIdPath, errRecordDatabaseIdRequired))
	} else if _, err := uuid.FromString(record.DatabaseID); err != nil {
		result = append(result, validation.Invalid(databaseIdPath, record.DatabaseID, errRecordInvalidDatabaseId))
	}

	if len(record.FormID) == 0 {
		result = append(result, validation.Required(formIdPath, errRecordFormIdRequired))
	} else if _, err := uuid.FromString(record.FormID); err != nil {
		result = append(result, validation.Invalid(formIdPath, record.FormID, errRecordInvalidFormId))
	}

	if _, ok := form.(types.SubFormInterface); ok {
		if record.OwnerID == nil || len(*record.OwnerID) == 0 {
			result = append(result, validation.Required(ownerIdPath, errRecordOwnerIdRequired))
		} else if _, err := uuid.FromString(*record.OwnerID); err != nil {
			result = append(result, validation.Invalid(ownerIdPath, *record.OwnerID, errRecordInvalidOwnerID))
		}
	}

	result = append(result, ValidateRecordValues(valuesPath, record.Values, form)...)

	return result
}

func ValidateRecordValues(path *validation.Path, recordValues types.FieldValues, form types.FormInterface) validation.ErrorList {
	var result validation.ErrorList
	if recordValues == nil {
		return append(result, validation.Required(path, errRecordValuesRequired))
	}

	// Keep track of what fields were sent as values
	recordFieldIDs := sets.NewString()
	recordFieldIndexes := map[string]int{}
	recordValueMap := map[string]types.FieldValue{}
	for i, recordValue := range recordValues {
		recordValueFieldID := recordValue.FieldID
		if recordFieldIDs.Has(recordValueFieldID) {
			result = append(result, validation.Duplicate(path.Index(i).Child("fieldId"), recordValueFieldID))
		} else {
			recordValueMap[recordValueFieldID] = recordValue
			recordFieldIndexes[recordValueFieldID] = i
			recordFieldIDs.Insert(recordValueFieldID)
		}
	}

	// Keep a map of the expected fields for that form
	expectedFieldMap := make(map[string]*types.FieldDefinition)
	// Also keep a list of expected field ids
	expectedFieldIDs := sets.NewString()

	for _, formField := range form.GetFields() {
		formFieldKind, err := formField.FieldType.GetFieldKind()
		if err != nil {
			return append(result, validation.InternalError(path, err))
		}
		// populate the expectedFieldMap and expectedFieldIDs
		if _, ok := supportedRecordFieldKindMap[formFieldKind]; ok {
			expectedFieldMap[formField.ID] = formField
			expectedFieldIDs.Insert(formField.ID)
		}
	}

	// checking for values that the user sent that should not have been provided
	extraneousFieldIDs := recordFieldIDs.Difference(expectedFieldIDs)
	if !extraneousFieldIDs.IsEmpty() {
		for _, extraneousFieldID := range extraneousFieldIDs.List() {
			extraneousFieldValueIndex := recordFieldIndexes[extraneousFieldID]
			fieldIdPath := path.Index(extraneousFieldValueIndex).Child("fieldId")
			result = append(result, validation.NotSupported(fieldIdPath, extraneousFieldID, expectedFieldIDs.List()))
		}
	}

	// checking for required fields that the user did not send
	missingFieldIDs := expectedFieldIDs.Difference(recordFieldIDs)
	for _, missingFieldID := range missingFieldIDs.List() {
		expectedField := expectedFieldMap[missingFieldID]
		if expectedField.Required {
			result = append(result, validation.Required(path, errFieldValueRequired))
		}
	}

	// validate the field values that the user provided
	for _, fieldID := range expectedFieldIDs.Intersection(recordFieldIDs).List() {
		recordFieldIndex := recordFieldIndexes[fieldID]
		fieldPath := path.Index(recordFieldIndex)
		recordFieldValue := recordValueMap[fieldID].Value
		expectedField := expectedFieldMap[fieldID]
		result = append(result, ValidateRecordValue(fieldPath, recordFieldValue, expectedField)...)
	}

	return result
}

func ValidateRecordValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {

	var result validation.ErrorList
	fieldKind, _ := field.FieldType.GetFieldKind()

	switch fieldKind {
	case types.FieldKindText:
		return ValidateRecordStringValue(path, value, field)
	case types.FieldKindReference:
		return ValidateRecordReferenceValue(path, value, field)
	case types.FieldKindMultilineText:
		return ValidateRecordStringValue(path, value, field)
	case types.FieldKindDate:
		return ValidateRecordDateValue(path, value, field)
	case types.FieldKindQuantity:
		return ValidateRecordQuantityValue(path, value, field)
	case types.FieldKindMonth:
		return ValidateRecordMonthValue(path, value, field)
	case types.FieldKindWeek:
		return ValidateRecordWeekValue(path, value, field)
	case types.FieldKindSingleSelect:
		return ValidateRecordSingleSelectValue(path, value, field)
	case types.FieldKindMultiSelect:
		return ValidateRecordMultiSelectValue(path, value, field)
	case types.FieldKindBoolean:
		return ValidateRecordBooleanValue(path, value, field)
	}
	return result
}

func ValidateRecordStringValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	_, result, done := getStringValue(path, value, field, validation.ErrorList{})
	if done {
		return result
	}
	return result
}

func ValidateRecordDateValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	stringValue, result, done := getStringValue(path, value, field, validation.ErrorList{})
	if done {
		return result
	}
	_, err := time.Parse("2006-01-02", stringValue)
	if err != nil {
		valuePath := path.Child("value")
		return append(result, validation.Invalid(valuePath, value, errRecordInvalidDate))
	}
	return result
}

func ValidateRecordMonthValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	stringValue, result, done := getStringValue(path, value, field, validation.ErrorList{})
	if done {
		return result
	}
	_, err := time.Parse("2006-01", stringValue)
	if err != nil {
		valuePath := path.Child("value")
		return append(result, validation.Invalid(valuePath, value, errRecordInvalidMonth))
	}
	return result
}

func ValidateRecordWeekValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	stringValue, result, done := getStringValue(path, value, field, validation.ErrorList{})
	if done {
		return result
	}
	valuePath := path.Child("value")

	_, err := dates.ParseIsoWeekTime(stringValue)

	if err != nil {
		return append(result, validation.Invalid(valuePath, value, errRecordInvalidWeek))
	}
	return result
}

func ValidateRecordQuantityValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	var result validation.ErrorList
	valuePath := path.Child("value")

	if value.Kind == types.NullValue {
		if field.Required {
			result = append(result, validation.Required(valuePath, errFieldValueRequired))
		}
		return result
	}

	if value.Kind != types.StringValue {
		result = append(result, validation.Invalid(valuePath, value.GetValue(), errFieldValueMustBeString))
		return result
	}

	_, err := strconv.Atoi(value.StringValue)
	if err != nil {
		return append(result, validation.Invalid(valuePath, value, errRecordInvalidQuantity))
	}

	// we don't assert the zero value for an int field
	return result
}

func ValidateRecordReferenceValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	stringValue, result, done := getStringValue(path, value, field, validation.ErrorList{})
	if done {
		return result
	}
	valuePath := path.Child("value")
	if _, err := uuid.FromString(stringValue); err != nil {
		return append(result, validation.Invalid(valuePath, value, errRecordInvalidReferenceUid))
	}
	return result
}

func ValidateRecordSingleSelectValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	stringValue, result, done := getStringValue(path, value, field, validation.ErrorList{})
	if done {
		return result
	}
	valuePath := path.Child("value")

	acceptedOptionIDs := sets.NewString()
	for _, option := range field.FieldType.SingleSelect.Options {
		acceptedOptionIDs.Insert(option.ID)
	}

	if !acceptedOptionIDs.Has(stringValue) {
		result = append(result, validation.NotSupported(valuePath, stringValue, acceptedOptionIDs.List()))
	}

	return result
}

func ValidateRecordMultiSelectValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	var result validation.ErrorList
	valuePath := path.Child("value")

	if value.Kind == types.ArrayValue {
		if len(value.ArrayValue) == 0 {
			if field.Required {
				return append(result, validation.Required(valuePath, errFieldValueRequired))
			} else {
				return result
			}
		}
	} else if value.Kind == types.NullValue {
		if field.Required {
			return append(result, validation.Required(valuePath, errFieldValueRequired))
		} else {
			return result
		}
	} else {
		return append(result, validation.Invalid(valuePath, value.GetValue(), errFieldValueMustBeArray))
	}

	acceptedOptionIDs := sets.NewString()
	for _, option := range field.FieldType.MultiSelect.Options {
		acceptedOptionIDs.Insert(option.ID)
	}

	seenValues := sets.NewString()
	for _, selectedOption := range value.ArrayValue {
		if !acceptedOptionIDs.Has(selectedOption) {
			result = append(result, validation.NotSupported(valuePath, selectedOption, acceptedOptionIDs.List()))
		}
		if seenValues.Has(selectedOption) {
			result = append(result, validation.Duplicate(valuePath, selectedOption))
		}
		seenValues.Insert(selectedOption)
	}

	return result
}

func ValidateRecordBooleanValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition) validation.ErrorList {
	stringValue, result, done := getStringValue(path, value, field, validation.ErrorList{})
	if done {
		return result
	}
	valuePath := path.Child("value")

	if stringValue != "true" && stringValue != "false" {
		result = append(result, validation.NotSupported(valuePath, stringValue, []string{"true", "false"}))
	}

	return result
}

func getStringValue(path *validation.Path, value types.StringOrArray, field *types.FieldDefinition, result validation.ErrorList) (string, validation.ErrorList, bool) {
	valuePath := path.Child("value")
	if value.Kind == types.NullValue {
		if field.Required {
			result = append(result, validation.Required(valuePath, errFieldValueRequired))
		}
		return "", result, true
	}

	if value.Kind != types.StringValue {
		result = append(result, validation.Invalid(valuePath, value.GetValue(), errFieldValueMustBeString))
		return "", result, true
	}

	if field.Required && strings.TrimSpace(value.StringValue) == "" {
		result = append(result, validation.Required(valuePath, errFieldValueRequired))
		return "", result, true
	}

	return value.StringValue, result, false
}
