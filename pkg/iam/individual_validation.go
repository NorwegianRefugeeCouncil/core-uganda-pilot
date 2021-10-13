package iam

import (
	"fmt"
	"github.com/nrc-no/core/pkg/validation"
)

func ValidateIndividual(individual *Individual, path *validation.Path) validation.ErrorList {
	// Validate Party
	errs := ValidateParty(individual.Party, path)
	if !individual.Party.HasPartyType(IndividualPartyType.ID) {
		errs = append(errs, validation.InternalError(path.Child("partyTypeIds"), fmt.Errorf("individual is missing the individual party type")))
	}

	return errs
}
