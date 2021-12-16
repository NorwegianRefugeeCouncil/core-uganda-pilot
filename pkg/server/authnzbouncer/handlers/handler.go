package handlers

import (
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/server/authnzbouncer/authenticators"
	"github.com/nrc-no/core/pkg/server/authnzbouncer/authorizers"
	"github.com/nrc-no/core/pkg/utils"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"sync"
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

		// we do not write directly to the response headers.
		// if we run authorization and authentication in parallel,
		// those might add headers to the response even if
		// one of the goroutines fail.
		// To prevent leaking any authorization/authentication
		// in the case of a failure, we delay writing to the
		// headers until both the authorization and authentication
		// are successful
		respHeaders := url.Values{}

		// We also ensure a thread-safe writing to the response headers
		// using a mutex. The authorization part does not write any
		// headers yet, but the implementation is there for when it does.
		respLock := sync.Mutex{}

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

			respLock.Lock()
			defer respLock.Unlock()
			respHeaders.Set("X-Remote-Subject", authenticationResponse.Sub)
			respHeaders.Set("X-Remote-Username", authenticationResponse.PreferredUsername)

			return nil
		})

		// Waiting for both authorization and authentication to finish.
		if err := g.Wait(); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		for key, values := range respHeaders {
			w.Header()[key] = values
		}
		w.WriteHeader(http.StatusOK)
	}
}
