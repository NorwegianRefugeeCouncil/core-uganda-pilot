package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postParty(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var party Party
	if err := s.bind(req, &party); err != nil {
		s.error(w, err)
		return
	}

	p := &party

	if p.ID == "" {
		p.ID = uuid.NewV4().String()
	}

	errList := ValidateParty(p, validation.NewPath(""))
	if len(errList) > 0 {
		status := validation.Status{
			Status:  validation.Failure,
			Code:    http.StatusUnprocessableEntity,
			Message: "invalid party",
			Errors:  errList,
		}
		s.json(w, status.Code, status)
		return
	}

	if err := s.partyStore.create(ctx, p); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, p)
}
