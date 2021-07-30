package attachments

import (
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostAttachment(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload Attachment
	if err := s.Bind(req, &payload); err != nil {
		s.Error(w, err)
		return
	}

	att := &payload

	errList := ValidateAttachment(att, &validation.Path{})
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid attachment")
		s.Error(w, &status)
		return
	}

	att.ID = uuid.NewV4().String()

	if err := s.Store.Create(ctx, att); err != nil {
		s.Error(w, err)
		return
	}

	s.json(w, http.StatusOK, att)
}
