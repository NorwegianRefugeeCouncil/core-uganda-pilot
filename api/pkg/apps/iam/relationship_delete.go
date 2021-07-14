package iam

import (
	"net/http"
)

func (s *Server) deleteRelationship(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	err := s.relationshipStore.delete(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}
}
