package iam

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (s *Server) postIndividual(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var individual Individual
	if err := s.bind(req, &individual); err != nil {
		s.error(w, err)
		return
	}

	if individual.ID == "" {
		individual.ID = uuid.NewV4().String()
	}

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

	if err := s.individualStore.create(ctx, &individual); err != nil {
		s.error(w, err)
		return
	}

	s.json(w, http.StatusOK, individual)
}
