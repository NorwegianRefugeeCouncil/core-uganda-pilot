package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var party Party
	if err := s.Bind(req, &party); err != nil {
		s.Error(w, err)
		return
	}

	p := &party
	p.ID = uuid.NewV4().String()

	if err := s.PartyStore.Create(ctx, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, p)
}
