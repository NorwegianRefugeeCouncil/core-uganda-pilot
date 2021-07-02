package cms

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostCase(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Case
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	kase := &payload
	kase.ID = uuid.NewV4().String()
	subject := ctx.Value("Subject").(string)
	kase.CreatorID = subject

	if err := s.caseStore.Create(ctx, kase); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, kase)
}
