package iam

import (
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

	if err := s.partyStore.update(ctx, r); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, r)
}
