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
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
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
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, updated)

}
