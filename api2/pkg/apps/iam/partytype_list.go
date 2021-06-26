package iam

import (
	"net/http"
)

func (s *Server) ListPartyTypes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &PartyTypeListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	ret, err := s.PartyTypeStore.List(ctx, *listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
