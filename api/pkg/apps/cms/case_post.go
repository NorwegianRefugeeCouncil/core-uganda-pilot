package cms

import (
	"github.com/nrc-no/core/pkg/validation"
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
	if kase.Template != nil {
		errList := ValidateCase(kase, &validation.Path{})
		if len(errList) > 0 {
			status := errList.Status(http.StatusUnprocessableEntity, "invalid case")
			s.Error(w, &status)
			return
		}
	} else {
		ct, err := s.caseTypeStore.Get(ctx, kase.CaseTypeID)
		if err != nil {
			s.Error(w, err)
			return
		}
		kase.Template = ct.Template
	}

	kase.ID = uuid.NewV4().String()

	subject, ok := ctx.Value("Subject").(string)
	if ok {
		kase.CreatorID = subject
	}
	if err := s.caseStore.Create(ctx, kase); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, kase)
}
