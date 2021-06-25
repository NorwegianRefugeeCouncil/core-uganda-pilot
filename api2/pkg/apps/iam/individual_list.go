package iam

import (
	"net/http"
)

func (s *Server) ListIndividuals(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	listOptions := IndividualListOptions{
		PartyTypeIDs: req.URL.Query()["partyTypeIds"],
	}

	list, err := s.IndividualStore.List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, list)
}
