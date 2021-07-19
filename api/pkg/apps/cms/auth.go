package cms

import (
	"context"
	"github.com/nrc-no/core/pkg/auth"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
)

func (s *Server) WithAuth() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			if s.environment == "Development" && auth.SetDevAuthenticatedUserSubject(handler, w, req) {
				return
			}

			token, err := auth.AuthHeaderTokenSource(req).GetToken()
			if err != nil {
				s.Error(w, err)
				return
			}

			res, err := s.HydraAdmin.IntrospectOAuth2Token(&admin.IntrospectOAuth2TokenParams{
				Token:      token,
				Context:    req.Context(),
				HTTPClient: s.HydraHttpClient,
			})
			if err != nil {
				s.Error(w, err)
				return
			}

			if !*res.Payload.Active {
				s.Error(w, err)
				return
			}
			ctx := req.Context()
			ctx = context.WithValue(ctx, "Subject", res.Payload.Sub)
			req = req.WithContext(ctx)
			handler.ServeHTTP(w, req)
		})
	}
}
