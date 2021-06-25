package iam

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ListParties(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	listOptions := &PartyListOptions{
		PartyTypeID: req.URL.Query().Get("partyTypeId"),
	}

	ret, err := s.PartyStore.List(ctx, *listOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
