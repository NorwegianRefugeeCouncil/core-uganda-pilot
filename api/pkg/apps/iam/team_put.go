package iam

import (
	"net/http"
)

func (s *Server) PutTeam(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var payload Team
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	r, err := s.TeamStore.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ID = id
	r.Name = payload.Name

	if err := s.TeamStore.Update(ctx, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, r)

}
