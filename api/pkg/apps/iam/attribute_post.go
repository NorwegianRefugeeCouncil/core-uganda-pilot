package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postAttributes(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var a Attribute
	if err := s.bind(req, &a); err != nil {
		s.error(w, err)
	}

	a.ID = uuid.NewV4().String()

	if err := s.attributeStore.create(ctx, &a); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, a)

}
