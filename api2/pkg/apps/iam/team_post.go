package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostTeam(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Team
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	team := &payload
	team.ID = uuid.NewV4().String()

	if err := s.TeamStore.Create(ctx, team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, team)
}
