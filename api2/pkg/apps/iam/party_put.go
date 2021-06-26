package iam

import (
	"net/http"
)

func (s *Server) PutParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	r, err := s.PartyStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	var payload Party
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	r.Attributes = payload.Attributes
	r.PartyTypeIDs = payload.PartyTypeIDs

	if err := s.PartyStore.Update(ctx, r); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, r)
}
