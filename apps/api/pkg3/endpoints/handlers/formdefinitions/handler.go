package formdefinitions

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/apps/api/pkg3/apis"
	"github.com/nrc-no/core/apps/api/pkg3/storage"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	resourceName = "formdefinitions"
	groupName    = "nrc.no"
	apiVersion   = "nrc.no/v1"
)

type Handler struct {
	storage storage.Interface
}

func (h *Handler) Get(w http.ResponseWriter, req *http.Request) {

	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "unable to get resource id", http.StatusBadRequest)
		return
	}

	ctx := req.Context()

	var formDefinition apis.FormDefinition
	if err := h.storage.Get(ctx, fmt.Sprintf("%s-%s-%s", groupName, resourceName, id), &formDefinition); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bodyBytes, err := json.Marshal(formDefinition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)

}

func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var formDefinition apis.FormDefinition
	if err := json.Unmarshal(bodyBytes, &formDefinition); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var out apis.FormDefinition
	if err := h.storage.Update(ctx, "core.nrc.no/formdefinitions/"+formDefinition.ObjectMeta.UID, &formDefinition, &out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Post(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var formDefinition apis.FormDefinition
	if err := json.Unmarshal(bodyBytes, &formDefinition); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formDefinition.ObjectMeta.UID = ""
	formDefinition.ObjectMeta.CreationTimestamp = time.Now().UTC()
	formDefinition.ObjectMeta.DeletionTimestamp = nil
	formDefinition.SetAPIVersion(apiVersion)

	var out apis.FormDefinition
	if err := h.storage.Create(ctx, "formdefinitions-"+formDefinition.ObjectMeta.UID, &formDefinition, &out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(&out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

}
