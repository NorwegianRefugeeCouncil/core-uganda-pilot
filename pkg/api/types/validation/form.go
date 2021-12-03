package validation

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"regexp"
	"strings"
)

const (
	formDefinitionMinNameLength              = 3
	formDefinitionMaxNameLength              = 64
	formNameRequired                         = "form name is required"
	formNameNoLeadingTrailingWhitespaces     = "form name cannot have leading or trailing whitespaces"
	formDatabaseIdRequired                   = "database id is required"
	uuidInvalid                              = "invalid uuid"
	formFieldsRequired                       = "must have at least 1 field"
	fieldNameMinLength                       = 2
	fieldNameMaxLength                       = 64
	errFieldNameRequired                     = "field name is required"
	errFieldNameInvalid                      = "field name is invalid. It can only contain alphanumeric characters and spaces"
	fieldCodeMinLength                       = 1
	fieldCodeMaxLength                       = 32
	errInvalidFieldCode                      = "field code can only contain alphanumeric characters. It also can only start with a letter"
	errFieldNameNoLeadingTrailingWhitespaces = "field name cannot have leading or trailing whitespaces"
	errKeyFieldMustBeRequired                = "field marked as key must be have required = true"
	errFieldTypesMultipleF                   = "only one field type can be specified but received %v"
	errOneFieldTypeRequired                  = "at least one field type must be specified"
	referenceFieldDatabaseIdRequired         = "database id is required"
	referenceFieldDatabaseIdInvalid          = "invalid database id"
	referenceFieldFormIdRequired             = "form id is required"
	referenceFieldFormIdInvalid              = "invalid form id"
)

func ValidateForm(form *types.FormDefinition) validation.ErrorList {
	var result validation.ErrorList
	result = append(result, validateFormNameFn(form.Name, validation.NewPath("name"))...)
	result = append(result, validateDatabaseIdFn(form.DatabaseID, validation.NewPath("databaseId"))...)
	result = append(result, validateFolderIdFn(form.FolderID, validation.NewPath("folderId"))...)
	result = append(result, ValidateFormFields(form.Fields, validation.NewPath("fields"))...)
	return result
}

type validateStrFn func(str string, path *validation.Path) validation.ErrorList

var (
	validateFormNameFn   validateStrFn
	validateDatabaseIdFn validateStrFn
	validateFolderIdFn   validateStrFn
)

func init() {
	validateFormNameFn = ValidateFormName
	validateDatabaseIdFn = ValidateDatabaseName
	validateFolderIdFn = ValidateFormFolderID
}

func ValidateFormName(formName string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList
	if len(formName) == 0 {
		result = append(result, validation.Required(path, formNameRequired))
		return result
	}
	if strings.TrimSpace(formName) != formName {
		result = append(result, validation.Invalid(path, formName, formNameNoLeadingTrailingWhitespaces))
	}
	if len(formName) < formDefinitionMinNameLength {
		result = append(result, validation.TooShort(path, formName, formDefinitionMinNameLength))
	}
	if len(formName) > formDefinitionMaxNameLength {
		result = append(result, validation.TooLong(path, formName, formDefinitionMaxNameLength))
	}
	return result
}

func ValidateFormDatabaseID(databaseID string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList
	if len(databaseID) == 0 {
		result = append(result, validation.Required(path, formDatabaseIdRequired))
		return result
	}
	if _, err := uuid.FromString(databaseID); err != nil {
		result = append(result, validation.Invalid(path, databaseID, uuidInvalid))
		return result
	}
	return result
}

func ValidateFormFolderID(folderId string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList
	if len(folderId) == 0 {
		return result
	}
	if _, err := uuid.FromString(folderId); err != nil {
		result = append(result, validation.Invalid(path, folderId, uuidInvalid))
		return result
	}
	return result
}

func ValidateFormFields(fields types.FieldDefinitions, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList
	if len(fields) == 0 {
		result = append(result, validation.Required(path, formFieldsRequired))
		return result
	}
	for i, field := range fields {
		result = append(result, ValidateFieldDefinition(field, path.Index(i))...)
	}
	return result
}

func ValidateFieldDefinition(field *types.FieldDefinition, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList
	result = append(result, ValidateFieldName(field.Name, path.Child("name"))...)
	result = append(result, ValidateFieldCode(field.Code, path.Child("code"))...)

	// Fields with key = true must have required = true
	if field.Key && !field.Required {
		result = append(result, validation.Invalid(path.Child("required"), field.Required, errKeyFieldMustBeRequired))
	}

	result = append(result, ValidateFieldType(field.FieldType, path.Child("fieldType"))...)

	return result
}

var fieldNameRegex = regexp.MustCompile("^[a-zA-Z0-9]+( [a-zA-Z0-9]+)*$")

func ValidateFieldName(fieldName string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList
	if len(fieldName) == 0 {
		result = append(result, validation.Required(path, errFieldNameRequired))
		return result
	}
	if len(fieldName) < fieldNameMinLength {
		result = append(result, validation.TooShort(path, fieldName, fieldNameMinLength))
	}
	if len(fieldName) > fieldNameMaxLength {
		result = append(result, validation.TooLong(path, fieldName, fieldNameMaxLength))
		// stop processing here in case we have a crazy long field name
		return result
	}
	if strings.TrimSpace(fieldName) != fieldName {
		result = append(result, validation.Invalid(path, fieldName, errFieldNameNoLeadingTrailingWhitespaces))
	}
	if !fieldNameRegex.MatchString(fieldName) {
		result = append(result, validation.Invalid(path, fieldName, errFieldNameInvalid))
	}
	return result
}

var fieldCodeRegex = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z0-9]*$")

func ValidateFieldCode(fieldCode string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList
	// code is optional
	if len(fieldCode) == 0 {
		return result
	}
	if len(fieldCode) < fieldCodeMinLength {
		result = append(result, validation.TooShort(path, fieldCode, fieldCodeMinLength))
	}
	if len(fieldCode) > fieldCodeMaxLength {
		result = append(result, validation.TooLong(path, fieldCode, fieldCodeMaxLength))
	}
	if !fieldCodeRegex.MatchString(fieldCode) {
		result = append(result, validation.Invalid(path, fieldCode, errInvalidFieldCode))
	}
	return result
}

func ValidateFieldType(fieldType types.FieldType, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	type val struct {
		name  string
		value interface{}
	}
	candidates := []val{
		{name: "text", value: fieldType.Text},
		{name: "subForm", value: fieldType.SubForm},
		{name: "date", value: fieldType.Date},
		{name: "reference", value: fieldType.Reference},
		{name: "month", value: fieldType.Month},
		{name: "multilineText", value: fieldType.MultilineText},
		{name: "quantity", value: fieldType.Quantity},
		{name: "singleSelect", value: fieldType.SingleSelect},
	}
	var vals []val
	var valNames []string
	for _, candidate := range candidates {
		if !reflect.ValueOf(candidate.value).IsNil() {
			vals = append(vals, candidate)
			valNames = append(valNames, candidate.name)
		}
	}

	if len(vals) == 0 {
		result = append(result, validation.Invalid(path, fieldType, errOneFieldTypeRequired))
		return result
	}

	if len(vals) > 1 {
		result = append(result, validation.Invalid(path, fieldType, fmt.Sprintf(errFieldTypesMultipleF, valNames)))
		return result
	}

	if fieldType.Text != nil {
		result = append(result, ValidateFieldTypeText(fieldType.Text, path.Child("text"))...)
	}
	if fieldType.SubForm != nil {
		result = append(result, ValidateFieldTypeSubForm(fieldType.SubForm, path.Child("subForm"))...)
	}
	if fieldType.Date != nil {
		result = append(result, ValidateFieldTypeDate(fieldType.Date, path.Child("date"))...)
	}
	if fieldType.Reference != nil {
		result = append(result, ValidateFieldTypeReference(fieldType.Reference, path.Child("reference"))...)
	}
	if fieldType.Month != nil {
		result = append(result, ValidateFieldTypeMonth(fieldType.Month, path.Child("month"))...)
	}
	if fieldType.MultilineText != nil {
		result = append(result, ValidateFieldTypeMultilineText(fieldType.MultilineText, path.Child("multilineText"))...)
	}
	if fieldType.Quantity != nil {
		result = append(result, ValidateFieldTypeQuantity(fieldType.Quantity, path.Child("quantity"))...)
	}
	if fieldType.SingleSelect != nil {
		result = append(result, ValidateFieldTypeSingleSelect(fieldType.SingleSelect, path.Child("singleSelect"))...)
	}

	return result
}

func ValidateFieldTypeText(ftText *types.FieldTypeText, path *validation.Path) validation.ErrorList {
	// noop
	var result []*validation.Error
	return result
}

func ValidateFieldTypeSubForm(ftSF *types.FieldTypeSubForm, path *validation.Path) validation.ErrorList {
	var result []*validation.Error
	result = append(result, ValidateFormName(ftSF.Name, path.Child("name"))...)

	result = append(result, ValidateFormFields(ftSF.Fields, path.Child("fields"))...)
	return result
}

func ValidateFieldTypeDate(ftDate *types.FieldTypeDate, path *validation.Path) validation.ErrorList {
	// noop
	var result []*validation.Error
	return result
}

func ValidateFieldTypeReference(ftRef *types.FieldTypeReference, path *validation.Path) validation.ErrorList {
	var result []*validation.Error
	if len(ftRef.DatabaseID) == 0 {
		result = append(result, validation.Required(path, referenceFieldDatabaseIdRequired))
		return result
	} else if _, err := uuid.FromString(ftRef.DatabaseID); err != nil {
		result = append(result, validation.Required(path, referenceFieldDatabaseIdInvalid))
		return result
	}
	if len(ftRef.FormID) == 0 {
		result = append(result, validation.Required(path, referenceFieldFormIdRequired))
		return result
	} else if _, err := uuid.FromString(ftRef.FormID); err != nil {
		result = append(result, validation.Required(path, referenceFieldFormIdInvalid))
		return result
	}
	return result
}

func ValidateFieldTypeMonth(ftMonth *types.FieldTypeMonth, path *validation.Path) validation.ErrorList {
	// noop
	var result []*validation.Error
	return result
}

func ValidateFieldTypeMultilineText(ftMultilineText *types.FieldTypeMultilineText, path *validation.Path) validation.ErrorList {
	// noop
	var result []*validation.Error
	return result
}

func ValidateFieldTypeQuantity(ftQuantity *types.FieldTypeQuantity, path *validation.Path) validation.ErrorList {
	// noop
	var result []*validation.Error
	return result
}

func ValidateFieldTypeSingleSelect(ftSingleSelect *types.FieldTypeSingleSelect, path *validation.Path) validation.ErrorList {
	// TODO
	var result []*validation.Error
	return result
}
