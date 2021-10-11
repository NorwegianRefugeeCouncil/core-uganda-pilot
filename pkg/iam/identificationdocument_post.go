package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postIdentificationDocument(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var a IdentificationDocument
	if err := s.bind(req, &a); err != nil {
		s.error(w, err)
	}

	a.ID = uuid.NewV4().String()

	if err := s.identificationDocumentStore.create(ctx, &a); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, a)
}
