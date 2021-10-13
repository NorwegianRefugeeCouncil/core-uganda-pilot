package cms

import "net/http"

func (s *Server) ListCaseTypes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &CaseTypeListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	ret, err := s.caseTypeStore.list(ctx, *listOptions)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}
