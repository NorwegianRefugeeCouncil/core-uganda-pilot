package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var party Party
	if err := s.bind(req, &party); err != nil {
		s.error(w, err)
		return
	}

	p := &party

	p.ID = uuid.NewV4().String()

	if err := s.partyStore.create(ctx, p); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, p)
}
