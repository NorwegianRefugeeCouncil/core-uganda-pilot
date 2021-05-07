package formdefinitions

import "net/http"

// List formDefinitions
func (h *Handler) List(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	if req.Header.Get("Upgrade") == "websocket" {
		h.Watch(w, req)
		return
	}

	var list = h.scope.NewList()
	if err := h.storage.List(ctx, list); err != nil {
		h.scope.Error(err, w, req)
		return
	}

	transformResponseObject(ctx, h.scope, req, w, http.StatusOK, list)
}
