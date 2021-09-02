package cms

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

func (s *Server) PutCase(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var id string
	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var payload Case
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	// Get Case from store
	kase, err := s.caseStore.get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	// Update struct
	update := updateCaseStruct(kase, payload)

	// Perform validation
	errList := ValidateCase(update, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid case")
		s.Error(w, &status)
		return
	}

	// if no validation errors, assume the case is Done if intake case
	if kase.IntakeCase {
		kase.Done = true
	} else {
		kase.Done = payload.Done
	}

	if err := s.caseStore.Update(ctx, kase); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, kase)

}
