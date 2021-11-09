package handlers

import (
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/server/auth/authenticators"
	"github.com/nrc-no/core/pkg/server/auth/authorizers"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func HandleAuth(
	authenticator authenticators.Authenticator,
	authorizer authorizers.Authorizer,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		authenticationResponse, err := authenticator.Authenticate(req)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		w.Header().Set("X-Remote-Subject", authenticationResponse.Sub)
		w.Header().Set("X-Remote-Username", authenticationResponse.PreferredUsername)

		authorizationResponse, err := authorizer.Authorize(req)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		if !authorizationResponse.Active {
			utils.ErrorResponse(w, meta.NewUnauthorized("token is not active"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
