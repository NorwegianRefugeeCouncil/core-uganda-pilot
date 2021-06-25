package iam

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

func (s *Server) PostRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var payload RelationshipType
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p := &payload
	p.ID = uuid.NewV4().String()

	if err := s.RelationshipTypeStore.Create(ctx, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
