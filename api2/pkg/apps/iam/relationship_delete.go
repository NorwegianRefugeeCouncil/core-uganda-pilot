package iam

import (
	"net/http"
)

func (s *Server) DeleteRelationship(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	err := s.RelationshipStore.Delete(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}
}
