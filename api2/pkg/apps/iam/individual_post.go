package iam

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

func (s *Server) PostIndividual(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var individual Individual
	if err := json.Unmarshal(bodyBytes, &individual); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(individual)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
