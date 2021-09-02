package cms

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func (s *Server) PostComment(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var payload Comment
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	subject := ctx.Value("Subject").(string)

	payload.AuthorID = subject
	now := time.Now().UTC()
	payload.CreatedAt = now
	payload.UpdatedAt = now
	payload.ID = uuid.NewV4().String()

	if err := s.commentStore.Create(ctx, &payload); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, &payload)

}
