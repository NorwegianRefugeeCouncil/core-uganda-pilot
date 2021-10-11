package attachments

import "net/http"

func (s *Server) ListAttachments(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	listOptions := &AttachmentListOptions{}
	if err := listOptions.UnmarshalQueryParameters(req.URL.Query()); err != nil {
		s.Error(w, err)
		return
	}

	ret, err := s.store.List(ctx, *listOptions)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)

}
