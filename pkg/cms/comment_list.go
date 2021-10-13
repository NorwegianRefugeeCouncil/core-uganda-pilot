package cms

import (
	"net/http"
)

func (s *Server) ListComments(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var listOptions = &CommentListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.error(w, err)
		return
	}

	comments, err := s.commentStore.List(ctx, *listOptions)
	if err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, comments)
}
