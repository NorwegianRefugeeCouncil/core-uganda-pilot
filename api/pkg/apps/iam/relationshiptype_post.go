package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload RelationshipType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	p := &payload
	p.ID = uuid.NewV4().String()

	if err := s.relationshipTypeStore.Create(ctx, p); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, p)
}
