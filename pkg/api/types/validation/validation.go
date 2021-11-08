package validation

import (
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
	"regexp"
)

func ValidateDatabase(db *types.Database) validation.ErrorList {
	allErrs := validation.ErrorList{}
	allErrs = append(allErrs, ValidateDatabaseName(db.Name, validation.NewPath("name"))...)
	return allErrs
}

var databaseNameRegex = regexp.MustCompile("^[A-Za-z0-9-_]+(?: [A-Za-z0-9-_]+)*$")
var databaseNameMinLength = 3
var databaseNameMaxLength = 32
var invalidDatabaseNameMsg = "Invalid database name. Supported values are (a-z A-Z - _)"

func ValidateDatabaseName(name string, field *validation.Path) validation.ErrorList {
	allErrs := validation.ErrorList{}
	if len(name) < databaseNameMinLength {
		allErrs = append(allErrs, validation.TooShort(field, name, databaseNameMinLength))
	}
	if len(name) > databaseNameMaxLength {
		allErrs = append(allErrs, validation.TooLong(field, name, databaseNameMaxLength))
	} else if !databaseNameRegex.MatchString(name) {
		// don't run the regex if database name is too long in case we have a crazy long database name
		allErrs = append(allErrs, validation.Invalid(field, name, invalidDatabaseNameMsg))
	}
	return allErrs
}
