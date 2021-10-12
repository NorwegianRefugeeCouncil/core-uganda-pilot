package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

func (s *Server) putPartyAttributeDefinition(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	// Unmarshal request payload
	var payload *PartyAttributeDefinition
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	// Get PartyAttributeDefinition from store
	partyAttributeDefinition, err := s.partyAttributeDefinitionStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}
	// Update
	update := updatePartyAttributesDefinitionStruct(partyAttributeDefinition, payload)

	// Perform validation
	errList := ValidatePartyAttributeDefinition(update, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid attribute")
		s.error(w, &status)
		return
	}

	// Persist changes
	if err := s.partyAttributeDefinitionStore.update(ctx, update); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, update)
}

func updatePartyAttributesDefinitionStruct(pAD *PartyAttributeDefinition, update *PartyAttributeDefinition) *PartyAttributeDefinition {
	// make copy
	result := *pAD

	// update
	if len(update.CountryID) > 0 {
		result.CountryID = update.CountryID
	}
	if len(update.PartyTypeIDs) > 0 {
		result.PartyTypeIDs = update.PartyTypeIDs
	}
	result.FormControl = update.FormControl
	result.IsPersonallyIdentifiableInfo = update.IsPersonallyIdentifiableInfo

	return &result
}
