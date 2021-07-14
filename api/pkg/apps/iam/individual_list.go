package iam

import (
	"net/http"
)

func (s *Server) listIndividuals(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &IndividualListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	list, err := s.individualStore.list(ctx, *listOptions)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, list)
}
