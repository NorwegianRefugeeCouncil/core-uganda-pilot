package cms

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

func (s *Server) PutCase(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	// Unmarshal request payload
	var payload *Case
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	// Get Case from store
	kase, err := s.caseStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	// Update struct
	update := updateCaseStruct(kase, payload)

	// Perform validation
	errList := ValidateCase(update, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid case")
		s.error(w, &status)
		return
	}

	// Persist case changes
	if err := s.caseStore.update(ctx, update); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, update)
}

func updateCaseStruct(kase *Case, update *Case) *Case {
	// make copy
	result := *kase

	// update
	if update.FormData != nil {
		result.FormData = update.FormData
	}
	result.Done = update.Done

	return &result
}
