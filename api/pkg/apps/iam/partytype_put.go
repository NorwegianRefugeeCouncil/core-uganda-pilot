package iam

import (
	"net/http"
)

func (s *Server) putPartyType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	r, err := s.partyTypeStore.Get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	var payload PartyType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	r.Name = payload.Name
	r.IsBuiltIn = payload.IsBuiltIn

	if err := s.partyTypeStore.Update(ctx, r); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, r)
}
