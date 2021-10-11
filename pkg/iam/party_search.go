package iam

import (
	"net/http"
)

func (s *Server) searchParties(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var listOptions PartySearchOptions
	if err := s.bind(req, &listOptions); err != nil {
		s.error(w, err)
		return
	}

	ret, err := s.partyStore.list(ctx, PartySearchOptions{
		PartyIDs: listOptions.PartyIDs,
		PartyTypeIDs: listOptions.PartyTypeIDs,
		Attributes:   listOptions.Attributes,
		SearchParam:  listOptions.SearchParam,
	})
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}
