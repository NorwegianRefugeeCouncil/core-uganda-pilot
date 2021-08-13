package cms

import (
	"net/http"
)

func (s *Server) GetCase(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	ret, err := s.caseStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}
	if ret.Done {
		ret.Template = ret.Template.MarkAsReadonly()
	}
	s.JSON(w, http.StatusOK, ret)
}
