package iam

import (
	"net/http"
)

func (s *Server) ListAttributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	listOptions := AttributeListOptions{
		PartyTypeIDs: req.URL.Query()["partyTypeIds"],
	}

	list, err := s.AttributeStore.List(ctx, listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, list)
}
