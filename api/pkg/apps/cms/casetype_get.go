package cms

import (
	"net/http"
)

func (s *Server) GetCaseType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	ret, err := s.caseTypeStore.get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}
