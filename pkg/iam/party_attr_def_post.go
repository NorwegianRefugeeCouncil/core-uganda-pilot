package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postPartyAttributeDefinition(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload PartyAttributeDefinition
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	payload.ID = uuid.NewV4().String()
	if len(payload.FormControl.Name) == 0 {
		payload.FormControl.Name = payload.ID
	}

	errList := ValidatePartyAttributeDefinition(&payload, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid attribute")
		s.error(w, &status)
		return
	}

	if err := s.partyAttributeDefinitionStore.create(ctx, &payload); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, payload)
}
