package login

import (
	"fmt"
	"github.com/nrc-no/core/pkg/auth"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) WithAuth() func(handler http.Handler) http.Handler {

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			logrus.Trace("retrieving Authorization header")

			token, err := auth.AuthHeaderTokenSource(req).GetToken()
			if err != nil {
				s.Error(w, fmt.Errorf("failed to get token from Authorization header: %v", err))
				return
			}

			logrus.Trace("inspecting Oauth2 Token")

			res, err := s.HydraAdmin.IntrospectOAuth2Token(&admin.IntrospectOAuth2TokenParams{
				Token:      token,
				Context:    req.Context(),
				HTTPClient: nil,
			})
			if err != nil {
				s.Error(w, fmt.Errorf("failed to introspect token: %v", err))
				return
			}

			if !*res.Payload.Active {
				s.Error(w, fmt.Errorf("token is not active"))
				return
			}

			logrus.Trace("Authentication successful")

			handler.ServeHTTP(w, req)
		})
	}
}
