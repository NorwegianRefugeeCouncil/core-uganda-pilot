package iam

import (
	"net/http"
)

func (s *Server) putRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var payload RelationshipType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	r, err := s.relationshipTypeStore.Get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	r.FirstPartyRole = payload.FirstPartyRole
	r.SecondPartyRole = payload.SecondPartyRole
	r.Name = payload.Name
	r.Rules = payload.Rules
	r.IsDirectional = payload.IsDirectional

	if err := s.relationshipTypeStore.Update(ctx, r); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, r)
}
