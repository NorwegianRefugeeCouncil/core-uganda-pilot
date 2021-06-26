package iam

import (
	"net/http"
)

func (s *Server) ListIndividuals(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &IndividualListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	list, err := s.IndividualStore.List(ctx, *listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, list)
}
