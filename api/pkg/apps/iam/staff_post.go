package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostStaff(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Staff
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	staff := &payload
	staff.ID = uuid.NewV4().String()

	if err := s.StaffStore.Create(ctx, staff); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, staff)
}
