package iam

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ListRelationships(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	qry := req.URL.Query()

	listOptions := RelationshipListOptions{
		RelationshipTypeID: qry.Get("relationshipTypeId"),
		FirstPartyId:       qry.Get("firstPartyId"),
		SecondParty:        qry.Get("secondPartyId"),
		EitherParty:        qry.Get("eitherPartyId"),
	}

	ret, err := s.RelationshipStore.List(ctx, listOptions)
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
