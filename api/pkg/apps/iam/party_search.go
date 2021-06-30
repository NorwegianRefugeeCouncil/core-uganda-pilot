package iam

import (
	"net/http"
)

func (s *Server) SearchParties(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var listOptions PartySearchOptions
	if err := s.Bind(req, &listOptions); err != nil {
		s.Error(w, err)
		return
	}

	ret, err := s.PartyStore.List(ctx, PartySearchOptions{
		PartyTypeIDs: listOptions.PartyTypeIDs,
		Attributes:   listOptions.Attributes,
		SearchParam:  listOptions.SearchParam,
	})
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
