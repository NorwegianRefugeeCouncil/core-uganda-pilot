package cms

import "net/http"

func (s *Server) ListCases(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &CaseListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	ret, err := s.caseStore.List(ctx, *listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, ret)

}
