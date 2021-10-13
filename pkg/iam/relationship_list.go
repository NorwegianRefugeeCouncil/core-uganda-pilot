package iam

import (
	"net/http"
)

func (s *Server) listRelationships(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &RelationshipListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	ret, err := s.relationshipStore.list(ctx, *listOptions)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}