package cms

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

func (s *Server) putCaseType(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var payload *CaseType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	caseType, err := s.caseTypeStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	// Update
	update := updateCaseTypeStruct(caseType, payload)

	errList := ValidateCaseType(update, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid CaseType")
		s.error(w, &status)
		return
	}

	if err := s.caseTypeStore.update(ctx, update); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, update)

}

func updateCaseTypeStruct(caseType *CaseType, update *CaseType) *CaseType {
	// make copy
	result := *caseType

	if len(update.TeamID) > 0 {
		result.TeamID = update.TeamID
	}
	if len(update.PartyTypeID) > 0 {
		result.PartyTypeID = update.PartyTypeID
	}
	if len(update.Name) > 0 {
		result.Name = update.Name
	}
	result.Form = update.Form
	result.IntakeCaseType = update.IntakeCaseType
	return &result
}
