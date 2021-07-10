package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postRelationship(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Relationship
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	p := &payload

	if p.ID == "" {
		p.ID = uuid.NewV4().String()
	}

	if err := s.relationshipStore.create(ctx, p); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, p)
}
