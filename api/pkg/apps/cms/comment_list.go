package cms

import (
	"net/http"
)

func (s *Server) ListComments(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var listOptions = &CommentListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	comments, err := s.commentStore.List(ctx, *listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, comments)

}
