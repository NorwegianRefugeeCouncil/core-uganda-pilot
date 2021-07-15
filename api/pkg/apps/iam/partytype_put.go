package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	"net/http"
)

func (s *Server) putPartyType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var id string

	if !s.getPathParam("id", w, req, &id) {
		return
	}

	r, err := s.partyTypeStore.Get(ctx, id)
	if err != nil {
		s.error(w, err)
		return
	}

	var payload PartyType
	if err := s.bind(req, &payload); err != nil {
		s.error(w, err)
		return
	}

	r.Name = payload.Name
	r.IsBuiltIn = payload.IsBuiltIn

	errList := ValidatePartyType(r, validation.NewPath(""))
	if len(errList) > 0 {
		status := validation.Status{
			Status:  validation.Failure,
			Code:    http.StatusUnprocessableEntity,
			Message: "invalid PartyType",
			Errors:  errList,
		}
		s.json(w, status.Code, status)
		return
	}

	if err := s.partyTypeStore.Update(ctx, r); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, r)
}
