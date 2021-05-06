package formdefinitions

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/nrc-no/core/apps/api/pkg/apis"
	"github.com/nrc-no/core/apps/api/pkg/endpoints"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	"github.com/nrc-no/core/apps/api/pkg/watch"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	resourceName = "formdefinitions"
	groupName    = "core.nrc.no"
	apiVersion   = "core.nrc.no/v1"
)

type Handler struct {
	storage storage.Interface
}

func (h *Handler) Get(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var requestInfo = ctx.Value("requestInfo").(*endpoints.RequestInfo)
	key := strings.ToLower(path.Join(requestInfo.APIGroup, requestInfo.APIResource, requestInfo.ResourceID))

	var formDefinition apis.FormDefinition
	if err := h.storage.Get(ctx, key, &formDefinition); err != nil {
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

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var requestInfo = ctx.Value("requestInfo").(*endpoints.RequestInfo)
	key := strings.ToLower(path.Join(requestInfo.APIGroup, requestInfo.APIResource))

	if req.Header.Get("Upgrade") == "websocket" {
		h.Watch(w, req)
		return
	}

	var formDefinition apis.FormDefinition
	if err := h.storage.Get(ctx, key, &formDefinition); err != nil {
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

	var requestInfo = ctx.Value("requestInfo").(*endpoints.RequestInfo)
	key := strings.ToLower(path.Join(requestInfo.APIGroup, requestInfo.APIResource, requestInfo.ResourceID))

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
	if err := h.storage.Update(ctx, key, &formDefinition, &out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

}

func (h *Handler) Post(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var requestInfo = ctx.Value("requestInfo").(*endpoints.RequestInfo)
	key := strings.ToLower(path.Join(requestInfo.APIGroup, requestInfo.APIResource))

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
	if err := h.storage.Create(ctx, key, &formDefinition, &out); err != nil {
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

func (h *Handler) Watch(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var requestInfo = ctx.Value("requestInfo").(*endpoints.RequestInfo)
	key := strings.ToLower(path.Join(requestInfo.APIGroup, requestInfo.APIResource))

	u := websocket.Upgrader{}
	conn, err := u.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	err = h.storage.Watch(ctx, key, &apis.FormDefinition{}, func(eventType string, obj runtime.Object) {
		if err := conn.WriteJSON(&watch.Event{
			Type:   eventType,
			Object: obj,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
