package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

func (s *Server) putParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	r, err := s.partyStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	var payload Party
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	r.Attributes = payload.Attributes
	r.PartyTypeIDs = payload.PartyTypeIDs

	errList := ValidateParty(r, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid party")
		s.json(w, status.Code, status)
		return
	}

	if err := s.partyStore.update(ctx, r); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, r)
}
