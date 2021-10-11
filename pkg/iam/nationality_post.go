package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postNationality(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Nationality
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	p := &payload

	p.ID = uuid.NewV4().String()

	if err := s.nationalityStore.create(ctx, p); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, p)

}
