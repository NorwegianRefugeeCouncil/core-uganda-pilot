package iam

import (
	"net/http"
)

func (s *Server) PutIndividual(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var individual Individual
	if err := s.Bind(req, &individual); err != nil {
		s.Error(w, err)
		return
	}

	_, err := s.IndividualStore.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.IndividualStore.Upsert(ctx, &individual); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, &individual)

}
