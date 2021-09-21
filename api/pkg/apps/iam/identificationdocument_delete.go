package iam

import (
	"net/http"
)

func (s *Server) deleteIdentificationDocument(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	err := s.identificationDocumentStore.delete(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}
}
