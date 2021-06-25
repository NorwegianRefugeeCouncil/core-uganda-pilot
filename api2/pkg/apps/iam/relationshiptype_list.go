package iam

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ListRelationshipTypes(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	qry := req.URL.Query()

	listOptions := RelationshipTypeListOptions{
		PartyType: qry.Get("partyType"),
	}

	ret, err := s.RelationshipTypeStore.List(ctx, listOptions)
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
