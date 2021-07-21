package cms

import (
	"net/http"
)

func (s *Server) PutCase(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var id string
	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var payload Case
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	kase, err := s.caseStore.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	kase.Done = payload.Done
	kase.FormData = payload.FormData

	if err := s.caseStore.Update(ctx, kase); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, kase)

}
