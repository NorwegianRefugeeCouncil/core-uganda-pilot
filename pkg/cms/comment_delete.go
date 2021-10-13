package cms

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) DeleteComment(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := mux.Vars(req)["id"]
	if err := s.commentStore.Delete(ctx, id); err != nil {
		s.error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
