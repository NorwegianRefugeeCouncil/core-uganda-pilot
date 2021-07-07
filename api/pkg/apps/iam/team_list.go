package iam

import (
	"net/http"
)

func (s *Server) listTeams(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ret, err := s.teamStore.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.json(w, http.StatusOK, ret)
}
