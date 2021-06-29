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
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, comment)

}
