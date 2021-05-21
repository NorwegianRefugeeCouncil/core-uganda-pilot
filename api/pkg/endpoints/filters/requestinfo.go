package filters

import (
	"fmt"
	request2 "github.com/nrc-no/core/api/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
)

// WithRequestInfo retrieves and inserts the RequestInfo in the request context.context
func WithRequestInfo(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		info, err := request2.NewRequestInfo(req)
		if err != nil {
			responsewriters.InternalError(w, req, fmt.Errorf("failed to create RequestInfo: %v", err))
			return
		}
		req = req.WithContext(request2.WithRequestInfo(ctx, info))
		handler.ServeHTTP(w, req)
	})
}
