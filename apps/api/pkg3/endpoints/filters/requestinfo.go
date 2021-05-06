package filters

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg3/endpoints"
	"net/http"
)

func WithRequestInfo(handler http.Handler, requestInfoResolver endpoints.RequestInfoResolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		info, err := requestInfoResolver.NewRequestInfo(req)
		if err != nil {
			http.Error(w, "unable to create RequestInfo", http.StatusInternalServerError)
			return
		}
		req = req.WithContext(context.WithValue(ctx, "requestInfo", info))
		handler.ServeHTTP(w, req)
	})
}
