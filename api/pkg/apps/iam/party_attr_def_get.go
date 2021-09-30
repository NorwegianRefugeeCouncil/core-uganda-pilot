package iam

import (
	"net/http"
)

func (s *Server) getPartyAttributeDefinition(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	a, err := s.partyAttributeDefinitionStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, a)
}
