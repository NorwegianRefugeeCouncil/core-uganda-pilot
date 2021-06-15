package relationships

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store}
}

func (h *Handler) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("id not found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ret, err := h.store.Get(ctx, id)
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

type ListOptions struct {
	Party string
}

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	qry := req.URL.Query()

	listOptions := &ListOptions{
		Party: qry.Get("Party"),
	}

	ret, err := h.store.List(ctx, *listOptions)
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

func (h *Handler) Put(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("id not found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r, err := h.store.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var payload Relationship
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.EndOfRelationship = payload.EndOfRelationship
	r.StartOfRelationship = payload.StartOfRelationship
	r.FirstParty = payload.FirstParty
	r.SecondParty = payload.SecondParty

	if err := h.store.Update(ctx, r); err != nil {
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

func (h *Handler) Post(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var payload Relationship
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p := &payload
	p.ID = uuid.NewV4().String()

	if err := h.store.Create(ctx, p); err != nil {
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
