package formdefinitions

import (
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints"
	"io/ioutil"
	"net/http"
)

// Update formDefinition
func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var requestInfo = ctx.Value("requestInfo").(*endpoints.RequestInfo)

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

	var formDefinition v1.FormDefinition
	_, _, err = h.scope.Serializer.Decode(bodyBytes, &h.scope.Kind, &formDefinition)
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

	var out v1.FormDefinition
	if err := h.storage.Update(ctx, requestInfo.ResourceID, &formDefinition, &out); err != nil {
		h.scope.Error(err, w, req)
		return
	}

	transformResponseObject(ctx, h.scope, req, w, http.StatusOK, &out)

}
