package iam

import (
	"net/http"
)

func (s *Server) putIndividual(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var individual Individual
	if err := s.bind(req, &individual); err != nil {
		s.error(w, err)
		return
	}

	_, err := s.individualStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	if err := s.individualStore.upsert(ctx, &individual); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, &individual)

}
