package iam

import (
	"net/http"
)

func (s *Server) PutPartyType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	r, err := s.PartyTypeStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	var payload PartyType
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	r.Name = payload.Name
	r.IsBuiltIn = payload.IsBuiltIn

	if err := s.PartyTypeStore.Update(ctx, r); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, r)
}
