package iam

import (
	"github.com/nrc-no/core/pkg/validation"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postAttributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var a Attribute
	if err := s.bind(req, &a); err != nil {
		s.error(w, err)
	}

	if a.ID == "" {
		a.ID = uuid.NewV4().String()
	}

	errList := ValidateAttribute(&a, validation.NewPath(""))
	if len(errList) > 0 {
		status := errList.Status(http.StatusUnprocessableEntity, "invalid attribute")
		s.json(w, status.Code, status)
		return
	}

	if err := s.attributeStore.create(ctx, &a); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, a)

}
