package iam

import (
	"net/http"
)

func (s *Server) GetStaff(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	ret, err := s.StaffStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
