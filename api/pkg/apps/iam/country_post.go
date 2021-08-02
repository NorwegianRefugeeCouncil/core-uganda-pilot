package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postCountry(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Country
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	country := &payload
	country.ID = uuid.NewV4().String()

	if err := s.countryStore.Create(ctx, country); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, country)
}
