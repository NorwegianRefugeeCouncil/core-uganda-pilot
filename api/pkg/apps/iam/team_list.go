package iam

import (
	"net/http"
)

func (s *Server) ListTeams(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ret, err := s.TeamStore.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
