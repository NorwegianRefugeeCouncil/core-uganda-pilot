package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postPartyType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var payload PartyType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	p := &payload
	p.ID = uuid.NewV4().String()

	errList := ValidatePartyType(p, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid party type")
		s.error(w, &status)
		return
	}

	if err := s.partyTypeStore.Create(ctx, p); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, p)
}