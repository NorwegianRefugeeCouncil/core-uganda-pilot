package cms

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) GetComment(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	id := mux.Vars(req)["id"]

	comment, err := s.commentStore.Get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, comment)

}
