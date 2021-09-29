package form

import (
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/validation"
)

const MaxStringLength = 200
const MaxTextLength = 2000

// ValidateControlValue returns a validation.ErrorList containing validation errors encountered
// when evaluating a value for a given Control context.
func ValidateControlValue(control Control, value []string, path *validation.Path) validation.ErrorList {
	var err *validation.Error
	var result validation.ErrorList
	var childPath = path.Child(control.Name)

	// check required first
	if control.Validation.Required && utils.AllEmpty(value) {
		err = validation.Required(childPath, "this field is required")
		result = append(result, err)
		return result
	}

	// validate according to control type
	switch control.Type {
	case Text:
		result = validateTextControl(control, value, path)
	case Email:
		result = validateEmailControl(control, value, path)
	case Phone:
		result = validatePhoneControl(control, value, path)
	case URL:
		result = validateURLControl(control, value, path)
	case Date:
		result = validateDateControl(control, value, path)
	case Textarea:
		result = validateTextareaControl(control, value, path)
	case Dropdown:
		result = validateDropdownControl(control, value, path)
	case Boolean, Checkbox, Radio:
		result = validateCheckboxControl(control, value, path)
	case Taxonomy:
		result = validateTaxonomyControl(control, value, path)
	case File:
		result = validateFileControl(control, value, path)
	case Time:
		result = validateTimeControl(control, value, path)
	}
	return validation.ErrorList{}
}

func validateTextControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	result = validateSingleStringControl(value, path, result)

	return result
}

func validateSingleStringControl(value []string, path *validation.Path, result validation.ErrorList) validation.ErrorList {
	result = validateSingleSliceValue(value, path, result)
	text := value[0]
	result = validateMaxStringLength(text, path, result)
	return result
}

func validateMaxStringLength(text string, path *validation.Path, result validation.ErrorList) validation.ErrorList {
	// the text should not exceed a certain limit
	if len(text) > MaxStringLength {
		err := validation.TooLong(path, text, MaxStringLength)
		result = append(result, err)
	}
	return result
}

func validateSingleSliceValue(value []string, path *validation.Path, result validation.ErrorList) validation.ErrorList {
	// there should only be one item in the slice
	if len(value) > 1 {
		err := validation.TooMany(path, len(value), 1)
		result = append(result, err)
	}
	return result
}

func validateEmailControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	result = validateSingleStringControl(value, path, result)

	email := value[0]
	// the email should be valid
	if !validation.IsValidEmail(email) {
		err := validation.Invalid(path, email, validation.InvalidEmailDetail)
		result = append(result, err)
	}

	return result
}

func validatePhoneControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	result = validateSingleStringControl(value, path, result)

	phone := value[0]
	// the phone number should be valid
	if !validation.IsValidPhone(phone) {
		err := validation.Invalid(path, phone, validation.InvalidPhoneDetail)
		result = append(result, err)
	}

	return result
}

func validateURLControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	result = validateSingleStringControl(value, path, result)

	url := value[0]
	if !validation.IsValidURL(url) {
		err := validation.Invalid(path, url, validation.InvalidURLDetail)
		result = append(result, err)
	}

	return result
}

func validateDateControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	result = validateSingleStringControl(value, path, result)

	date := value[0]
	if !validation.IsValidDate(date) {
		err := validation.Invalid(path, date, validation.InvalidDateDetail)
		result = append(result, err)
	}

	return result
}

func validateTimeControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	result = validateSingleStringControl(value, path, result)

	date := value[0]
	if !validation.IsValidTime(date) {
		err := validation.Invalid(path, date, validation.InvalidTimeDetail)
		result = append(result, err)
	}

	return result
}

func validateTextareaControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	result = validateSingleSliceValue(value, path, result)

	text := value[0]
	// verify max length is not exceeded
	if len(text) > MaxTextLength {
		err := validation.TooLong(path, text, MaxTextLength)
		result = append(result, err)
	}
	return result
}

func validateDropdownControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	if !control.Multiple {
		result = validateSingleSliceValue(value, path, result)
	}

	return result
}

func validateCheckboxControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// check required options
	for i, option := range control.CheckboxOptions {
		if option.Required && !utils.Contains(value, option.Value) {
			optionPath := path.Index(i)
			err := validation.Required(optionPath, "this is required")
			result = append(result, err)
		}
	}

	return result
}

func validateTaxonomyControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// TODO implement me

	return result
}

func validateFileControl(control Control, value []string, path *validation.Path) validation.ErrorList {
	var result validation.ErrorList

	// TODO implement me

	return result
}
