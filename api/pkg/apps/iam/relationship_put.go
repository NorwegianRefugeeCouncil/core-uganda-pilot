package iam

import (
	"net/http"
)

func (s *Server) PutRelationship(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	r, err := s.RelationshipStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	var payload Relationship
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	r.ID = id
	r.RelationshipTypeID = payload.RelationshipTypeID
	r.FirstPartyID = payload.FirstPartyID
	r.SecondPartyID = payload.SecondPartyID

	if err := s.RelationshipStore.Update(ctx, r); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, r)
}
