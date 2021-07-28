package iam

import (
	"net/http"
)

func (s *Server) listCountrys(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ret, err := s.countryStore.List(ctx)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}
