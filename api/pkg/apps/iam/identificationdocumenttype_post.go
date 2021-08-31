package iam

import "net/http"

func (s *Server) postIdentificationDocumentType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &IdentificationDocumentTypeListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	list, err := s.identificationDocumetTypeStore.list(ctx, *listOptions)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, list)
}
