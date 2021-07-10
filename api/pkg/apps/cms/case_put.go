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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kase.Done = payload.Done
	kase.ParentID = payload.ParentID
	kase.TeamID = payload.TeamID

	if err := s.caseStore.Update(ctx, kase); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, kase)

}
