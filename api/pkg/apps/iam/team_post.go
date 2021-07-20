package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postTeam(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Team
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	team := &payload
	team.ID = uuid.NewV4().String()

	if err := s.teamStore.Create(ctx, team); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, team)
}
