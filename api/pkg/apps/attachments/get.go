package attachments

import "net/http"

func (s *Server) GetAttachment(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	ret, err := s.Store.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	s.json(w, http.StatusOK, ret)
}
