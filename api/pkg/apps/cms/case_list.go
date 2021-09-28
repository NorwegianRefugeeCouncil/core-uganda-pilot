package cms

import "net/http"

func (s *Server) ListCases(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &CaseListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	ret, err := s.caseStore.list(ctx, *listOptions)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)

}
