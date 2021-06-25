package iam

import (
	"net/http"
)

func (s *Server) ListParties(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	listOptions := &PartyListOptions{
		PartyTypeID: req.URL.Query().Get("partyTypeId"),
	}

	ret, err := s.PartyStore.List(ctx, *listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
