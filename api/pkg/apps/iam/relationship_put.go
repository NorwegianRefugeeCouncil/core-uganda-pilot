package iam

import (
	"net/http"
)

func (s *Server) putRelationship(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	r, err := s.relationshipStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	var payload Relationship
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	r.ID = id
	r.RelationshipTypeID = payload.RelationshipTypeID
	r.FirstPartyID = payload.FirstPartyID
	r.SecondPartyID = payload.SecondPartyID

	if err := s.relationshipStore.update(ctx, r); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, r)
}
