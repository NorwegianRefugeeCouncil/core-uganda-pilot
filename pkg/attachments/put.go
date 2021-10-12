package attachments

import "net/http"

func (s *Server) PutAttachment(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var id string
	if !s.GetPathParam("id", w, req, &id) {
		return
	}

	var payload Attachment
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	att, err := s.store.Get(ctx, id)
	if err != nil {
		s.Error(w, err)
		return
	}

	att.AttachedToID = payload.AttachedToID
	att.Body = payload.Body

	if err := s.store.Update(ctx, att); err != nil {
		s.Error(w, err)
		return
	}

	s.json(w, http.StatusOK, att)
}
