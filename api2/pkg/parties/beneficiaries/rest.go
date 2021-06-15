package beneficiaries

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
	return &Handler{
		store: store,
	}
}

func (h *Handler) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("id not found in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := h.store.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)

}

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	listOptions := ListOptions{
		PartyTypes: req.URL.Query()["partyTypes"],
	}

	list, err := h.store.List(ctx, listOptions)
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

func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var beneficiary Beneficiary
	if err := json.Unmarshal(bodyBytes, &beneficiary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	beneficiary.ID = uuid.NewV4().String()

	attrs := map[string][]string{}
	for key, values := range beneficiary.Attributes {
		if len(values) == 0 {
			continue
		}
		for _, value := range values {
			if len(value) == 0 {
				continue
			}
			attrs[key] = append(attrs[key], value)
		}
	}

	beneficiary.Attributes = attrs

	if err := h.store.Create(ctx, &beneficiary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(beneficiary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)

}

func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	id, ok := mux.Vars(req)["id"]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("no id in path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := h.store.Get(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var beneficiary Beneficiary
	if err := json.Unmarshal(bodyBytes, &beneficiary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.store.Upsert(ctx, &beneficiary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Get(w, req)

}
