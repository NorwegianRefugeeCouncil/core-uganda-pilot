package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostAttribute(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var a Attribute
	if err := s.Bind(req, &a); err != nil {
		s.Error(w, err)
	}

	a.ID = uuid.NewV4().String()

	if err := s.AttributeStore.Create(ctx, &a); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, a)

}
