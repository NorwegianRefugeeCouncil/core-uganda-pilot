package formdefinitions

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/writers"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"strings"
)

// Get a formDefinition
func (h *Handler) Get(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var requestInfo = ctx.Value("requestInfo").(*endpoints.RequestInfo)
	key := strings.ToLower(path.Join(requestInfo.APIGroup, requestInfo.APIResource, requestInfo.ResourceID))

	var formDefinition v1.FormDefinition
	if err := h.storage.Get(ctx, key, &formDefinition); err != nil {
		h.scope.Error(err, w, req)
		return
	}

	transformResponseObject(ctx, h.scope, req, w, http.StatusOK, &formDefinition)
}

func transformResponseObject(ctx context.Context, scope *handlers.RequestScope, req *http.Request, w http.ResponseWriter, statusCode int, result runtime.Object) {
	err := scope.Serializer.Encode(result, w)
	if err == nil {
		return
	}

	status := writers.ErrorToApiStatus(err)
	candidateStatusCode := int(status.Code)

	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		w.WriteHeader(candidateStatusCode)
	}

	output, err := runtime.Encode(scope.Serializer, status)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		output = []byte(fmt.Sprintf("%s: %s", status.Reason, status.Message))
	}
	if _, err := w.Write(output); err != nil {
		logrus.Error("unable to write fallback json response: %v", err)
	}
	return
}
