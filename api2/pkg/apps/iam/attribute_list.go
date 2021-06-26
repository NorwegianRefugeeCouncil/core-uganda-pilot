package iam

import (
	"net/http"
)

func (s *Server) ListAttributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	listOptions := &AttributeListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	list, err := s.AttributeStore.List(ctx, *listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, list)
}
