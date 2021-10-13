package iam

import (
	"net/http"
)

func (s *Server) listParties(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &PartyListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	options := &PartySearchOptions{
		Attributes:  listOptions.Attributes,
		SearchParam: listOptions.SearchParam,
	}
	if len(listOptions.PartyTypeID) > 0 {
		options.PartyTypeIDs = []string{listOptions.PartyTypeID}
	}

	ret, err := s.partyStore.list(ctx, *options)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}
