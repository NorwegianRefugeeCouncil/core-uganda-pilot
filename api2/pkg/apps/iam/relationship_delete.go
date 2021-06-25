package iam

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) DeleteRelationship(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("id not found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := s.RelationshipStore.Delete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
