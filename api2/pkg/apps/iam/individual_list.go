package iam

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ListIndividuals(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	listOptions := IndividualListOptions{
		PartyTypeIDs: req.URL.Query()["partyTypeIds"],
	}

	list, err := s.IndividualStore.List(ctx, listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
