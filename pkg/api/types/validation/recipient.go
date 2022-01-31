package validation

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
)

const (
	errRecipientMultipleKeyFields         = "Recipient forms can have at most one key field"
	errRecipientMustHaveReferenceKeyField = "The key field for recipient forms must be a Reference field"
)

// ValidateRecipientForm validates recipient forms
func ValidateRecipientForm(form *types.FormDefinition) validation.ErrorList {
	var result validation.ErrorList
	if form.Type != types.RecipientFormType {
		return result
	}

	fieldsPath := validation.NewPath("fields")

	keyField, keyFieldIndex, singleKeyFieldErrs := validateRecipientFormHasSingleKeyField(form.Fields, fieldsPath)
	result = append(result, singleKeyFieldErrs...)

	if keyField != nil {
		keyFieldRefErrs := validateRecipientKeyFieldIsReference(keyField, fieldsPath.Index(keyFieldIndex))
		result = append(result, keyFieldRefErrs...)
	}

	return result
}

func validateRecipientKeyFieldIsReference(
	keyField *types.FieldDefinition,
	path *validation.Path,
) validation.ErrorList {
	var result validation.ErrorList

	keyFieldKind, err := keyField.FieldType.GetFieldKind()
	if err != nil {
		result = append(result, validation.InternalError(path, err))
		return result
	}

	if keyFieldKind != types.FieldKindReference {
		result = append(result, validation.Invalid(path.Child("fieldType"), keyFieldKind, errRecipientMustHaveReferenceKeyField))
		return result
	}

	return result
}

// validateRecipientFormHasSingleKeyField validates that a recipient form has at most a single key field
func validateRecipientFormHasSingleKeyField(
	fields types.FieldDefinitions,
	path *validation.Path,
) (*types.FieldDefinition, int, validation.ErrorList) {
	var result validation.ErrorList
	var keyField *types.FieldDefinition
	var keyFieldIndex = -1
	for i, field := range fields {
		keyPath := path.Index(i).Child("key")
		if !field.Key {
			continue
		}
		if keyField != nil {
			result = append(result, validation.Invalid(keyPath, field.Key, errRecipientMultipleKeyFields))
			return nil, -1, result
		}
		keyFieldIndex = i
		keyField = field
	}
	return keyField, keyFieldIndex, result
}
