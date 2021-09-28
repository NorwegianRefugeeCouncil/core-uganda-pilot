package cms

import (
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postCaseType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload CaseType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	payload.ID = uuid.NewV4().String()

	errList := ValidateCaseType(&payload, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid case type")
		s.error(w, &status)
		return
	}

	if err := s.caseTypeStore.create(ctx, &payload); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, &payload)
}
