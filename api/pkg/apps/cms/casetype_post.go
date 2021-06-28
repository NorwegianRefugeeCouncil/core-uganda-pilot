package cms

import (
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
	caseType.ID = uuid.NewV4().String()

	if err := s.caseTypeStore.Create(ctx, caseType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.JSON(w, http.StatusOK, caseType)
}
