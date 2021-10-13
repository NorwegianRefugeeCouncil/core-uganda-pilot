package iam

import (
	"net/http"
)

func (s *Server) getParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	ret, err := s.partyStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}
