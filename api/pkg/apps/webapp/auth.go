package webapp

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
	"strings"
)

func (s *Server) WithAuth() func(handler http.Handler) http.Handler {

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			if s.environment == "Development" {
				authUserEmail := req.Header.Get("X-E2E-Authenticated-User-Email")
				if len(authUserEmail) != 0 {

					ctx := req.Context()

					authUsers, err := s.iamAdminClient.Parties().Search(ctx, iam.PartySearchOptions{
						Attributes: map[string]string{
							iam.EMailAttribute.ID: authUserEmail,
						},
					})
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					if len(authUsers.Items) == 0 {
						err := fmt.Errorf("user not found")
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					authUser := authUsers.Items[0]

					ctx = req.Context()
					ctx = context.WithValue(ctx, "Subject", authUser.ID)
					req = req.WithContext(ctx)

					handler.ServeHTTP(w, req)
					return

				}
			}

			// Do no perform auth when user tries to access login or callback
			// Because that is performing auth itself
			if req.URL.Path == "/login" || req.URL.Path == "/callback" {
				handler.ServeHTTP(w, req)
				return
			}

			// Get the access token from the session
			token := s.sessionManager.GetString(req.Context(), "access-token")

			if len(token) == 0 {
				authHeader := req.Header.Get("Authorization")
				authHeaderParts := strings.Split(authHeader, " ")
				if len(authHeaderParts) == 2 && authHeaderParts[0] == "Bearer" {
					token = authHeaderParts[1]
				}
			}

			// Redirect to login if no access token
			if len(token) == 0 {
				http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
				return
			}

			// Verify access token
			res, err := s.HydraAdmin.IntrospectOAuth2Token(&admin.IntrospectOAuth2TokenParams{
				Token:      token,
				Context:    req.Context(),
				HTTPClient: nil,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Check id token
			profile := s.sessionManager.Get(req.Context(), "profile")

			fmt.Printf("%#v", profile)

			if !*res.Payload.Active {
				http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			handler.ServeHTTP(w, req)
		})
	}

}
