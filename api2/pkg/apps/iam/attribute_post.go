package iam

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

func (s *Server) PostAttribute(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var a Attribute
	if err := json.Unmarshal(bodyBytes, &a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.ID = uuid.NewV4().String()

	if err := s.AttributeStore.Create(ctx, &a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)

}
