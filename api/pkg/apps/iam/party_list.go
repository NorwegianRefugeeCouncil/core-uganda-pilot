package iam

import (
	"net/http"
)

func (s *Server) ListParties(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &PartyListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	ret, err := s.PartyStore.List(ctx, *listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
