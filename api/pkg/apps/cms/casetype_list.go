package cms

import "net/http"

func (s *Server) ListCaseTypes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &CaseTypeListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	ret, err := s.caseTypeStore.List(ctx, *listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, ret)

}
