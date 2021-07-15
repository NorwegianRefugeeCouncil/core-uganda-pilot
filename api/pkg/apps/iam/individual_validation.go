package iam

import "github.com/nrc-no/core/pkg/validation"

func ValidateIndividual(individual *Individual, path *validation.Path) validation.ErrorList {
	// Validate Party
	errs := ValidateParty(individual.Party, path)

	return errs
}
