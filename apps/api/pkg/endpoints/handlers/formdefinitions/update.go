package formdefinitions

import (
	"encoding/json"
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

// Update formDefinition
func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var requestInfo = ctx.Value("requestInfo").(*endpoints.RequestInfo)
	key := strings.ToLower(path.Join(requestInfo.APIGroup, requestInfo.APIResource, requestInfo.ResourceID))

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

	var formDefinition v1.FormDefinition
	if err := json.Unmarshal(bodyBytes, &formDefinition); err != nil {
		h.scope.Error(err, w, req)
		return
	}

	var out v1.FormDefinition
	if err := h.storage.Update(ctx, key, &formDefinition, &out); err != nil {
		h.scope.Error(err, w, req)
		return
	}

	responseBytes, err := json.Marshal(out)
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

}
