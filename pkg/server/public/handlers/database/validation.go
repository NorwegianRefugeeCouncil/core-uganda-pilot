package database

import (
	"strings"

	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/validation"
)

func ValidDBNameLength(name string) bool {
	return len(name) >= 3 && len(name) <= 32
}

func DbNameHasNoLeadingOrTrailingWhitespace(name string) bool {
	trimmed := strings.TrimSpace(name)
	return name == trimmed
}

func ValidateDBStruct(db *types.Database) bool {
	returnVal := validation.IsValidAlpha(db.Name)
	returnVal = returnVal && ValidDBNameLength(db.Name)
	returnVal = returnVal && DbNameHasNoLeadingOrTrailingWhitespace(db.Name)
	return returnVal
}
