package iam

import (
	"net/http"
)

func (s *Server) ListStaff(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ret, err := s.StaffStore.List(ctx, StaffListOptions{})
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
