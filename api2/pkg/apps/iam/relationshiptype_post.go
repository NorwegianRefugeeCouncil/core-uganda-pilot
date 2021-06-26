package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload RelationshipType
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	p := &payload
	p.ID = uuid.NewV4().String()

	if err := s.RelationshipTypeStore.Create(ctx, p); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, p)
}
