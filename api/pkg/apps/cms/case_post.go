package cms

import (
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostCase(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Case
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	kase := &payload
	kase.ID = uuid.NewV4().String()

	errList := ValidateCase(kase, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid case")
		s.error(w, &status)
		return
	}

	subject, ok := ctx.Value("Subject").(string)
	if ok {
		kase.CreatorID = subject
	}
	if err := s.caseStore.create(ctx, kase); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, kase)
}
