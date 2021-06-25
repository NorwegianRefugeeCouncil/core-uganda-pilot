package iam

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ListPartyTypes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ret, err := s.PartyTypeStore.List(ctx)
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
