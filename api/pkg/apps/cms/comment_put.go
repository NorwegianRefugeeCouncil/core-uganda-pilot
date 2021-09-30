package cms

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func (s *Server) PutComment(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id := mux.Vars(req)["id"]

	var payload Comment
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	updated, err := s.commentStore.Update(ctx, id, func(oldComment *Comment) (*Comment, error) {
		now := time.Now().UTC()
		newComment := &*oldComment
		newComment.UpdatedAt = now
		newComment.Body = payload.ID
		return newComment, nil
	})
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, updated)

}
