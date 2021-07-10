package iam

import (
	"net/http"
)

func (s *Server) putTeam(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var payload Team
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	r, err := s.teamStore.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ID = id
	r.Name = payload.Name

	if err := s.teamStore.Update(ctx, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.json(w, http.StatusOK, r)

}
