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

	options := &PartySearchOptions{
		Attributes:  listOptions.Attributes,
		SearchParam: listOptions.SearchParam,
	}
	if len(listOptions.PartyTypeID) > 0 {
		options.PartyTypeIDs = []string{listOptions.PartyTypeID}
	}

	ret, err := s.PartyStore.List(ctx, *options)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)
}
