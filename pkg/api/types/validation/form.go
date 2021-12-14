package validation

import (
	"fmt"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/utils/sets"
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"regexp"
	"strings"
)

const (
	errDatabaseIdRequired                    = "Database id is required"
	errFieldNameInvalid                      = "Field name is invalid. It can only contain alphanumeric characters and spaces"
	errFieldNameNoLeadingTrailingWhitespaces = "Field name cannot have leading or trailing whitespaces"
	errFieldNameRequired                     = "Field name is required"
	errFieldTypesMultipleF                   = "Only one field type can be specified but received %v"
	errFieldsRequired                        = "Form must have at least 1 field"
	errFormNameRequired                      = "Form name is required"
	errFormNameWhitespace                    = "Form name cannot have leading or trailing whitespaces"
	errInvalidFieldCode                      = "Field code can only contain alphanumeric characters. It also can only start with a letter"
	errInvalidUUID                           = "Invalid uuid"
	errKeyFieldMustBeRequired                = "Field marked as key must be have required = true"
	errMultiLineTextFieldCannotBeKeyField    = "Multiline text fields cannot key marked as key fields"
	errMultiSelectCannotBeKeyField           = "Multi select fields cannot key marked as key field"
	errOneFieldTypeRequired                  = "At least one field type must be specified"
	errReferenceFieldDatabaseIdInvalid       = "Invalid database id"
	errReferenceFieldDatabaseIdRequired      = "Database id is required"
	errReferenceFieldFormIdInvalid           = "Invalid form id"
	errReferenceFieldFormIdRequired          = "Form id is required"
	errSelectOptionNameInvalid               = "Name of the select option is invalid"
	errSelectOptionNameRequired              = "Name of the option is required"
	errSelectOptionsRequired                 = "At least 1 option must be specified"
	errSubFormCannotBeKeyOrRequiredField     = "Sub form fields cannot key marked as key or required fields"
	fieldCodeMaxLength                       = 32
	fieldCodeMinLength                       = 1
	fieldNameMaxLength                       = 128
	fieldNameMinLength                       = 2
	formMaxFieldCount                        = 60
	formNameMaxLength                        = 128
	formNameMinLength                        = 3
	selectFieldMaxNameLength                 = 128
	selectFieldMaxOptions                    = 100
	// todo: add maximum total number of fields (form fields + subform fields)
)

var (
	fieldNameRegex  = regexp.MustCompile("^[a-zA-Z0-9]+( [a-zA-Z0-9]+)*$")
	fieldCodeRegex  = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z0-9]*$")
	optionNameRegex = regexp.MustCompile("^[a-zA-Z0-9]+( [a-zA-Z0-9]+)*$")
)

func ValidateForm(form *types.FormDefinition) validation.ErrorList {
	var result validation.ErrorList
	result = append(result, ValidateFormName(form.Name, validation.NewPath("name"))...)
	result = append(result, ValidateFormDatabaseID(form.DatabaseID, validation.NewPath("databaseId"))...)
	result = append(result, ValidateFormFolderID(form.FolderID, validation.NewPath("folderId"))...)
	result = append(result, ValidateFormFields(form.Fields, validation.NewPath("fields"))...)
	return result
}

func ValidateFormName(formName string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// validate that the form name is not empty
	if len(formName) == 0 {
		result = append(result, validation.Required(path, errFormNameRequired))
		return result
	}

	// validates that the form name does not contain surrounding whitespaces
	if strings.TrimSpace(formName) != formName {
		result = append(result, validation.Invalid(path, formName, errFormNameWhitespace))
	}

	// validates that the form name does not exceed the max size
	if len(formName) < formNameMinLength {
		result = append(result, validation.TooShort(path, formName, formNameMinLength))
	}

	// validates that the form name is long enough
	if len(formName) > formNameMaxLength {
		result = append(result, validation.TooLong(path, formName, formNameMaxLength))
	}

	return result
}

func ValidateFormDatabaseID(databaseID string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// validates that the databaseId is present
	if len(databaseID) == 0 {
		result = append(result, validation.Required(path, errDatabaseIdRequired))
		return result
	}

	// validates that the databaseId is a valid uui
	if _, err := uuid.FromString(databaseID); err != nil {
		result = append(result, validation.Invalid(path, databaseID, errInvalidUUID))
		return result
	}

	return result
}

func ValidateFormFolderID(folderId string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// the folderId is optional
	if len(folderId) == 0 {
		return result
	}

	// validates that the folderId is a valid uuid
	if _, err := uuid.FromString(folderId); err != nil {
		result = append(result, validation.Invalid(path, folderId, errInvalidUUID))
		return result
	}
	return result
}

func ValidateFormFields(fields types.FieldDefinitions, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// validates that the form has at least 1 field
	if len(fields) == 0 {
		result = append(result, validation.Required(path, errFieldsRequired))
		return result
	}

	if len(fields) > formMaxFieldCount {
		result = append(result, validation.TooMany(path, len(fields), formMaxFieldCount))
		return result
	}

	// validates each field
	for i, field := range fields {
		result = append(result, ValidateFieldDefinition(field, path.Index(i))...)
	}

	return result
}

func ValidateFieldDefinition(field *types.FieldDefinition, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	namePath := path.Child("name")
	codePath := path.Child("code")
	requiredPath := path.Child("required")
	keyPath := path.Child("key")
	fieldTypePath := path.Child("fieldType")

	// validates the field name
	result = append(result, ValidateFieldName(field.Name, namePath)...)

	// validates the field code
	result = append(result, ValidateFieldCode(field.Code, codePath)...)

	// key fields must be required as well
	if field.Key && !field.Required {
		result = append(result, validation.Invalid(requiredPath, field.Required, errKeyFieldMustBeRequired))
	}

	if field.Required || field.Key {
		// sub form fields cannot be marked as required/key fields
		// sub form fields do not have a corresponding column on the form's sql table
		// our implementation uses NOT NULL and UNIQUE constraint, for which we need an actual SQL Column
		if field.FieldType.SubForm != nil {
			result = append(result, validation.Invalid(requiredPath, field.Key, errSubFormCannotBeKeyOrRequiredField))
		}
	}

	if field.Key {
		// prevent putting a key on a multiline text field.
		// This field type is for user textual input, which would never really be unique anyways.
		if field.FieldType.MultilineText != nil {
			result = append(result, validation.Invalid(keyPath, field.Key, errMultiLineTextFieldCannotBeKeyField))
		}
		// multi select fields cannot be marked as required/key fields
		// multi select fields do not have a corresponding column on the form's sql table
		// our implementation uses NOT NULL and UNIQUE constraint, for which we need an actual SQL Column
		if field.FieldType.MultiSelect != nil {
			result = append(result, validation.Invalid(keyPath, field.Key, errMultiSelectCannotBeKeyField))
		}
	}

	// validates that the field name is a valid form name if the field is a subform field
	if field.FieldType.SubForm != nil {
		result = append(result, ValidateFormName(field.Name, namePath)...)
	}

	// validates the field type
	result = append(result, ValidateFieldType(field.FieldType, fieldTypePath)...)

	return result
}

func ValidateFieldName(fieldName string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// validates that the field name is not empty
	if len(fieldName) == 0 {
		result = append(result, validation.Required(path, errFieldNameRequired))
		return result
	}

	// validates that the field name is long enough
	if len(fieldName) < fieldNameMinLength {
		result = append(result, validation.TooShort(path, fieldName, fieldNameMinLength))
	}

	// validates that the field name is not too short
	if len(fieldName) > fieldNameMaxLength {
		result = append(result, validation.TooLong(path, fieldName, fieldNameMaxLength))
		// stop processing here in case we have a crazy long field name
		return result
	}

	// validates that the field name does not contain surrounding whitespaces
	if strings.TrimSpace(fieldName) != fieldName {
		result = append(result, validation.Invalid(path, fieldName, errFieldNameNoLeadingTrailingWhitespaces))
		return result
	}

	// validates that the field name matches the regex
	if !fieldNameRegex.MatchString(fieldName) {
		result = append(result, validation.Invalid(path, fieldName, errFieldNameInvalid))
	}

	return result
}

func ValidateFieldCode(fieldCode string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// field code is optional
	if len(fieldCode) == 0 {
		return result
	}

	// keeping this in the case that we change the fieldCodeMinLength, which is = 1
	//goland:noinspection GoBoolExpressions
	if fieldCodeMinLength > 1 && len(fieldCode) < fieldCodeMinLength {
		result = append(result, validation.TooShort(path, fieldCode, fieldCodeMinLength))
	}

	// validates that the field code is not too long
	if len(fieldCode) > fieldCodeMaxLength {
		result = append(result, validation.TooLong(path, fieldCode, fieldCodeMaxLength))
	}

	// validates that the field code matches the regex
	if !fieldCodeRegex.MatchString(fieldCode) {
		result = append(result, validation.Invalid(path, fieldCode, errInvalidFieldCode))
	}

	return result
}

func ValidateFieldType(fieldType types.FieldType, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	textPath := path.Child("text")
	multiLineTextPath := path.Child("multilineText")
	monthPath := path.Child("month")
	weekPath := path.Child("week")
	datePath := path.Child("date")
	quantityPath := path.Child("quantity")
	referencePath := path.Child("reference")
	subFormPath := path.Child("subForm")
	singleSelectPath := path.Child("singleSelect")
	multiSelect := path.Child("multiSelect")

	// finds what kind of field type is defined
	var found []types.FieldKind
	for _, kind := range types.GetAllFieldKinds() {
		field, err := fieldType.GetFieldType(kind)
		if err != nil {
			result = append(result, validation.InternalError(path, err))
			return result
		}
		value := reflect.ValueOf(field)
		if value.Kind() == reflect.Ptr && !value.IsNil() {
			found = append(found, kind)
		}
	}

	// validates that the field type must be defined
	if len(found) == 0 {
		result = append(result, validation.Required(path, errOneFieldTypeRequired))
		return result
	}

	// validates that there is not more than 1 field type defined
	if len(found) > 1 {
		result = append(result, validation.TooLong(path, fmt.Sprintf(errFieldTypesMultipleF, found), 1))
		return result
	}

	fieldKind := found[0]

	// validates the field type
	switch fieldKind {
	case types.FieldKindText:
		result = append(result, ValidateFieldTypeText(fieldType.Text, textPath)...)
	case types.FieldKindMultilineText:
		result = append(result, ValidateFieldTypeMultilineText(fieldType.MultilineText, multiLineTextPath)...)
	case types.FieldKindMonth:
		result = append(result, ValidateFieldTypeMonth(fieldType.Month, monthPath)...)
	case types.FieldKindWeek:
		result = append(result, ValidateFieldTypeWeek(fieldType.Week, weekPath)...)
	case types.FieldKindDate:
		result = append(result, ValidateFieldTypeDate(fieldType.Date, datePath)...)
	case types.FieldKindQuantity:
		result = append(result, ValidateFieldTypeQuantity(fieldType.Quantity, quantityPath)...)
	case types.FieldKindReference:
		result = append(result, ValidateFieldTypeReference(fieldType.Reference, referencePath)...)
	case types.FieldKindSubForm:
		result = append(result, ValidateFieldTypeSubForm(fieldType.SubForm, subFormPath)...)
	case types.FieldKindSingleSelect:
		result = append(result, ValidateFieldTypeSingleSelect(fieldType.SingleSelect, singleSelectPath)...)
	case types.FieldKindMultiSelect:
		result = append(result, ValidateFieldTypeMultiSelect(fieldType.MultiSelect, multiSelect)...)
	default:
		result = append(result, validation.InternalError(path, fmt.Errorf("unknown field kind %v", fieldKind)))
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

	// validates the sub form fields
	fieldsPath := path.Child("fields")
	result = append(result, ValidateFormFields(ftSF.Fields, fieldsPath)...)

	return result
}

func ValidateFieldTypeDate(ftDate *types.FieldTypeDate, path *validation.Path) validation.ErrorList {
	// noop
	var result []*validation.Error
	return result
}

func ValidateFieldTypeReference(ftRef *types.FieldTypeReference, path *validation.Path) validation.ErrorList {
	var result []*validation.Error
	databaseIdPath := path.Child("databaseId")
	formIdPath := path.Child("formId")

	if len(ftRef.DatabaseID) == 0 {
		// validates that the database id is present
		result = append(result, validation.Required(databaseIdPath, errReferenceFieldDatabaseIdRequired))
		return result
	} else if _, err := uuid.FromString(ftRef.DatabaseID); err != nil {
		// validates that the database id is a well-formed uuid
		result = append(result, validation.Invalid(databaseIdPath, ftRef.DatabaseID, errReferenceFieldDatabaseIdInvalid))
		return result
	}

	if len(ftRef.FormID) == 0 {
		// validates that the form id is present
		result = append(result, validation.Required(formIdPath, errReferenceFieldFormIdRequired))
		return result
	} else if _, err := uuid.FromString(ftRef.FormID); err != nil {
		// validates that the form id is a well-formed uuid
		result = append(result, validation.Invalid(formIdPath, ftRef.FormID, errReferenceFieldFormIdInvalid))
		return result
	}
	return result
}

func ValidateFieldTypeMonth(ftMonth *types.FieldTypeMonth, path *validation.Path) validation.ErrorList {
	// noop
	var result []*validation.Error
	return result
}

func ValidateFieldTypeWeek(ftWeek *types.FieldTypeWeek, path *validation.Path) validation.ErrorList {
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

// ValidateFieldTypeSingleSelect validate a FieldTypeSingleSelect field
func ValidateFieldTypeSingleSelect(ftSingleSelect *types.FieldTypeSingleSelect, path *validation.Path) validation.ErrorList {
	return ValidateSelectOptions(ftSingleSelect.Options, path)
}

// ValidateFieldTypeMultiSelect validate a types.FieldTypeMultiSelect field
func ValidateFieldTypeMultiSelect(ftMultiSelect *types.FieldTypeMultiSelect, path *validation.Path) validation.ErrorList {
	return ValidateSelectOptions(ftMultiSelect.Options, path)
}

func ValidateSelectOptions(options []*types.SelectOption, path *validation.Path) validation.ErrorList {
	var result []*validation.Error

	optionsCount := len(options)
	optionsPath := path.Child("options")
	if optionsCount == 0 {
		// checking that we have at least one option
		result = append(result, validation.Required(optionsPath, errSelectOptionsRequired))
	} else if optionsCount > selectFieldMaxOptions {
		// checking that we don't exceed the max option count
		// we also don't process this further since this might be a crazy long list
		return append(result, validation.TooMany(optionsPath, optionsCount, selectFieldMaxOptions))
	}

	seenNames := sets.NewString()
	for i, option := range options {

		optionName := option.Name
		optionPath := optionsPath.Index(i)
		namePath := optionPath.Child("name")

		if len(optionName) == 0 {
			// option name is empty
			result = append(result, validation.Required(namePath, errSelectOptionNameRequired))
		} else if len(optionName) > selectFieldMaxNameLength {
			// option name is too long
			result = append(result, validation.TooLong(namePath, optionName, selectFieldMaxNameLength))
		} else if !optionNameRegex.MatchString(optionName) {
			// option name does not match regex
			result = append(result, validation.Invalid(namePath, optionName, errSelectOptionNameInvalid))
		} else if seenNames.Has(optionName) {
			// option name was seen in another option for the same field
			result = append(result, validation.Duplicate(namePath, optionName))
		} else {
			// valid name
			seenNames.Insert(optionName)
		}

	}

	return result
}
