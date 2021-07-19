package iam

import (
	"net/http"
)

func (s *Server) listAttributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	listOptions := &AttributeListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	list, err := s.attributeStore.list(ctx, *listOptions)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, list)
}
