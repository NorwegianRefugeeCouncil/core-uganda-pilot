package cms

import (
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostCaseType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload CaseType
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	caseType := &payload

	errList := ValidateCaseType(caseType, &validation.Path{})
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid case type")
		s.JSON(w, status.Code, status)
		return
	}

	caseType.ID = uuid.NewV4().String()

	if err := s.caseTypeStore.Create(ctx, caseType); err != nil {
		s.Error(w, err)
		return
	}
	s.JSON(w, http.StatusOK, caseType)
}
