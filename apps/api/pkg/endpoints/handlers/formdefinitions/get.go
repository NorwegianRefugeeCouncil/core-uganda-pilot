package formdefinitions

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/writers"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/request"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/serializer/versioning"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"strings"
)

// Get a formDefinition
func (h *Handler) Get(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	var requestInfo = ctx.Value("requestInfo").(*request.RequestInfo)
	key := strings.ToLower(path.Join(requestInfo.APIGroup, requestInfo.APIResource, requestInfo.ResourceID))

	obj, err := h.getter.Get(ctx, key, &metav1.GetOptions{})
	if err != nil {
		h.scope.Error(err, w, req)
		return
	}

	transformResponseObject(ctx, h.scope, req, w, http.StatusOK, obj)
}

func transformResponseObject(ctx context.Context, scope *handlers.RequestScope, req *http.Request, w http.ResponseWriter, statusCode int, result runtime.Object) {
	codec := versioning.NewDefaultingCodecForScheme(scope.Scheme, scope.Serializer, nil, scope.Kind.GroupVersion(), nil)
	err := codec.Encode(result, w)
	if err == nil {
		return
	}

	status := writers.ErrorToApiStatus(err)
	candidateStatusCode := int(status.Code)

	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		w.WriteHeader(candidateStatusCode)
	}

	output, err := runtime.Encode(codec, status)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		output = []byte(fmt.Sprintf("%s: %s", status.Reason, status.Message))
	}
	if _, err := w.Write(output); err != nil {
		logrus.Error("unable to write fallback json response: %v", err)
	}
	return
}

// WriteRawJSON writes a non-API object in JSON.
func WriteRawJSON(statusCode int, object interface{}, w http.ResponseWriter) {
	output, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(output)
}
