package handlers

import (
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/server/authnzbouncer/authenticators"
	"github.com/nrc-no/core/pkg/server/authnzbouncer/authorizers"
	"github.com/nrc-no/core/pkg/utils"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func HandleAuth(
	authenticator authenticators.Authenticator,
	authorizer authorizers.Authorizer,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// do not authorize preflight requests
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// perform authorization and authentication in parallel.
		// This is not a very time-consuming operation, though
		// this endpoint will be called to authorize/authenticate
		// all requests coming into core.
		// (This server runs as an `ext_authz` filter on envoy)
		// See https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter
		// Saving a few ms here and there on all requests will add up.
		g, _ := errgroup.WithContext(req.Context())

		// Authorizing the request
		g.Go(func() error {
			authorizationResponse, err := authorizer.Authorize(req)
			if err != nil {
				return err
			}
			if !authorizationResponse.Active {
				return meta.NewUnauthorized("token is not active")
			}
			return nil
		})

		// Authenticating the request
		g.Go(func() error {
			authenticationResponse, err := authenticator.Authenticate(req)
			if err != nil {
				return err
			}
			w.Header().Set("X-Remote-Subject", authenticationResponse.Sub)
			w.Header().Set("X-Remote-Username", authenticationResponse.PreferredUsername)
			return nil
		})

		if err := g.Wait(); err != nil {
			utils.ErrorResponse(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
