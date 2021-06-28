package iam

import (
	"net/http"
)

func (s *Server) PutRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var payload RelationshipType
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	r, err := s.RelationshipTypeStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	r.FirstPartyRole = payload.FirstPartyRole
	r.SecondPartyRole = payload.SecondPartyRole
	r.Name = payload.Name
	r.Rules = payload.Rules
	r.IsDirectional = payload.IsDirectional

	if err := s.RelationshipTypeStore.Update(ctx, r); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, r)
}
