package iam

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func (s *Server) PutRelationshipType(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("id not found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r, err := s.RelationshipTypeStore.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	r.FirstPartyRole = payload.FirstPartyRole
	r.SecondPartyRole = payload.SecondPartyRole
	r.Name = payload.Name
	r.Rules = payload.Rules
	r.IsDirectional = payload.IsDirectional

	if err := s.RelationshipTypeStore.Update(ctx, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
