package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostMembership(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Membership
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	p := &payload
	p.ID = uuid.NewV4().String()

	if err := s.MembershipStore.Create(ctx, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, p)

}
