package validation

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

const (
	errRecordInvalidDatabaseId  = "Invalid database ID"
	errRecordDatabaseIdRequired = "Database ID is required"
	errRecordInvalidFormId      = "Invalid form ID"
	errRecordFormIdRequired     = "Record form ID is required"
	errRecordOwnerIdRequired    = "Record owner ID is required"
	errRecordInvalidOwnerID     = "Record owner ID is invalid"
	errRecordValuesRequired     = "Record values are required"
	errInvalidFieldValueTypeF   = "Invalid value type for field. Expected %T, got %T"
	errRecordInvalidDate        = "Invalid date. Expected YYYY-mm-DD"
	errRecordInvalidMonth       = "Invalid date. Expected YYYY-mm"
	errRecordInvalidQuantity    = "Invalid quantity"
	errFieldValueRequired       = "Field value is required"
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
	types.FieldKindSingleSelect,
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

	if form.HasOwner() {
		if record.OwnerID == nil || len(*record.OwnerID) == 0 {
			result = append(result, validation.Required(ownerIdPath, errRecordOwnerIdRequired))
		} else if _, err := uuid.FromString(*record.OwnerID); err != nil {
			result = append(result, validation.Invalid(ownerIdPath, *record.OwnerID, errRecordInvalidOwnerID))
		}
	}

	result = append(result, ValidateRecordValues(valuesPath, record.Values, form)...)

	return result
}

func ValidateRecordValues(path *validation.Path, recordValues map[string]interface{}, form types.FormInterface) validation.ErrorList {
	var result validation.ErrorList
	if recordValues == nil {
		return append(result, validation.Required(path, errRecordValuesRequired))
	}

	// Keep track of what fields were sent as values
	recordFieldIDs := sets.NewString()
	for recordFieldID := range recordValues {
		recordFieldIDs.Insert(recordFieldID)
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
	extraneousFields := recordFieldIDs.Difference(expectedFieldIDs)
	if !extraneousFields.IsEmpty() {
		for _, extraneousItem := range extraneousFields.List() {
			valuePath := path.Key(extraneousItem)
			result = append(result, validation.NotSupported(valuePath, extraneousItem, expectedFieldIDs.List()))
		}
	}

	// checking for required fields that the user did not send
	missingFieldIDs := expectedFieldIDs.Difference(recordFieldIDs)
	for _, missingFieldID := range missingFieldIDs.List() {
		fieldPath := path.Key(missingFieldID)
		expectedField := expectedFieldMap[missingFieldID]
		if expectedField.Required {
			result = append(result, validation.Required(fieldPath, errFieldValueRequired))
		}
	}

	// validate the field values that the user provided
	for _, fieldKey := range expectedFieldIDs.Intersection(recordFieldIDs).List() {
		fieldPath := path.Key(fieldKey)
		recordFieldValue := recordValues[fieldKey]
		expectedField := expectedFieldMap[fieldKey]
		result = append(result, ValidateRecordValue(fieldPath, recordFieldValue, expectedField)...)
	}

	return result
}

func ValidateRecordValue(path *validation.Path, value interface{}, field *types.FieldDefinition) validation.ErrorList {

	var result validation.ErrorList
	fieldKind, _ := field.FieldType.GetFieldKind()

	found := false
	for _, kind := range supportedRecordFieldKinds {
		if kind == fieldKind {
			found = true
		}
	}
	if !found {
		return append(result, validation.NotSupported(path, fieldKind, supportedFieldKindNames))
	}

	switch fieldKind {
	case types.FieldKindText:
		result = append(result, ValidateRecordStringValue(path, value, field)...)
	case types.FieldKindReference:
		// TODO
	case types.FieldKindMultilineText:
		result = append(result, ValidateRecordStringValue(path, value, field)...)
	case types.FieldKindDate:
		result = append(result, ValidateRecordDateValue(path, value, field)...)
	case types.FieldKindQuantity:
		result = append(result, ValidateRecordQuantityValue(path, value, field)...)
	case types.FieldKindMonth:
		result = append(result, ValidateRecordMonthValue(path, value, field)...)
	case types.FieldKindSingleSelect:
	}
	return result
}

func ValidateRecordStringValue(path *validation.Path, value interface{}, field *types.FieldDefinition) validation.ErrorList {
	var result validation.ErrorList
	stringValue, ok := value.(string)
	if !ok {
		result = append(result, validation.Invalid(path, value, fmt.Sprintf(errInvalidFieldValueTypeF, "", value)))
		return result
	}
	if field.Required && strings.TrimSpace(stringValue) == "" {
		result = append(result, validation.Required(path, errFieldValueRequired))
	}
	return result
}

func ValidateRecordDateValue(path *validation.Path, value interface{}, field *types.FieldDefinition) validation.ErrorList {
	var result validation.ErrorList
	stringValue, ok := value.(string)
	if !ok {
		return append(result, validation.Invalid(path, value, fmt.Sprintf(errInvalidFieldValueTypeF, "", value)))
	}
	if field.Required && stringValue == "" {
		return append(result, validation.Required(path, errFieldValueRequired))
	}
	_, err := time.Parse("2006-01-02", stringValue)
	if err != nil {
		return append(result, validation.Invalid(path, value, errRecordInvalidDate))
	}
	return result
}

func ValidateRecordMonthValue(path *validation.Path, value interface{}, field *types.FieldDefinition) validation.ErrorList {
	var result validation.ErrorList
	stringValue, ok := value.(string)
	if !ok {
		return append(result, validation.Invalid(path, value, fmt.Sprintf(errInvalidFieldValueTypeF, "", value)))
	}
	if field.Required && stringValue == "" {
		return append(result, validation.Required(path, errFieldValueRequired))
	}
	_, err := time.Parse("2006-01", stringValue)
	if err != nil {
		return append(result, validation.Invalid(path, value, errRecordInvalidMonth))
	}
	return result
}

func ValidateRecordQuantityValue(path *validation.Path, value interface{}, field *types.FieldDefinition) validation.ErrorList {
	var result validation.ErrorList
	_, ok := value.(int)
	if !ok {
		return append(result, validation.Invalid(path, value, errRecordInvalidQuantity))
	}
	// we don't assert the zero value for an int field
	return result
}
