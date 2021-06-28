package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) PostIndividual(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var individual Individual
	if err := s.Bind(req, &individual); err != nil {
		s.Error(w, err)
		return
	}

	individual.ID = uuid.NewV4().String()

	attrs := map[string][]string{}
	for key, values := range individual.Attributes {
		if len(values) == 0 {
			continue
		}
		for _, value := range values {
			if len(value) == 0 {
				continue
			}
			attrs[key] = append(attrs[key], value)
		}
	}

	individual.Attributes = attrs

	if err := s.IndividualStore.Create(ctx, &individual); err != nil {
		s.Error(w, err)
		return
	}

	s.JSON(w, http.StatusOK, individual)
}
