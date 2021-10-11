package iam

import (
	"net/http"
)

func (s *Server) putCountry(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	var payload Country
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	r, err := s.countryStore.Get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	r.ID = id
	r.Name = payload.Name

	if err := s.countryStore.Update(ctx, r); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, r)

}
