package webapp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (s *Server) WithAuth() func(handler http.Handler) http.Handler {

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			// When in development mode, the user may provide a custom header
			// containing the email of the user s.he wants to impersonate,
			// effectively bypassing authorization
			if s.environment == "Development" {
				authUserEmail := req.Header.Get("X-E2E-Authenticated-User-Email")
				if len(authUserEmail) != 0 {
					s.dangerouslySetAuthenticatedUserUsingEmail(w, req, authUserEmail, handler)
					return
				}
			}

			// Do no perform auth when user tries to access login or callback
			// Because that is performing auth itself
			// If we would return a unauthorized response, then the user could
			// never log in, as logging in would itself be "unauthorized"
			if req.URL.Path == "/login" || req.URL.Path == "/callback" {
				handler.ServeHTTP(w, req)
				return
			}

			// Get the access token from the session
			session, err := s.sessionManager.Get(req)
			if err != nil {
				logrus.WithError(err).Errorf("failed to get session, attempting to clear session and redirect to login")

				// if the session manager returns a non-nil session with the error
				// try to use it to clear the session
				if session != nil {
					session.Options.MaxAge = -1
					if err = sessions.Save(req, w); err != nil {
						logrus.WithError(err).Errorf("failed to clear session!")
						s.Error(w, err)
						return
					}
				}

				// make a new state variable for hydra login flow
				b := make([]byte, 32)
				_, err := rand.Read(b)
				if err != nil {
					logrus.WithError(err).Errorf("failed to create a new state variable for hydra login flow!")
					s.Error(w, err)
					return
				}
				state := base64.StdEncoding.EncodeToString(b)

				// create a hydra login flow redirect url with the new state
				// variable, and redirect the user
				redirectUrl := s.publicOauth2Config.AuthCodeURL(state)
				http.Redirect(w, req, redirectUrl, http.StatusTemporaryRedirect)
				return
			}

			var token string

			accessTokenIntf, ok := session.Values["access-token"]
			if ok {
				accessTokenStr, ok := accessTokenIntf.(string)
				if ok {
					token = accessTokenStr
				}
			}

			// Retrieve token from Authorization header if the token
			// was not present in the session
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
				HTTPClient: s.HydraHTTPClient,
			})
			if err != nil {
				s.Error(w, err)
				return
			}

			// If the user is not active, then redirect to login
			if !*res.Payload.Active {
				http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
				return
			}

			handler.ServeHTTP(w, req)
		})
	}

}

func (s *Server) dangerouslySetAuthenticatedUserUsingEmail(w http.ResponseWriter, req *http.Request, authUserEmail string, handler http.Handler) {
	ctx := req.Context()

	authUsers, err := s.iamAdminClient.Parties().Search(ctx, iam.PartySearchOptions{
		Attributes: map[string]string{
			iam.EMailAttribute.ID: authUserEmail,
		},
	})
	if err != nil {
		s.Error(w, err)
		return
	}
	if len(authUsers.Items) == 0 {
		err := fmt.Errorf("user not found")
		s.Error(w, err)
		return
	}
	authUser := authUsers.Items[0]

	var profile = &Claims{}
	profile.Email = authUser.Get(iam.EMailAttribute.ID)
	profile.FamilyName = authUser.Get(iam.LastNameAttribute.ID)
	profile.GivenName = authUser.Get(iam.FirstNameAttribute.ID)
	profile.Subject = authUser.ID

	session, err := s.sessionManager.Get(req)
	if err != nil {
		s.Error(w, err)
		return
	}
	session.Values["profile"] = profile
	if err := session.Save(req, w); err != nil {
		s.Error(w, err)
		return
	}

	ctx = req.Context()
	ctx = context.WithValue(ctx, "Subject", authUser.ID)
	req = req.WithContext(ctx)

	handler.ServeHTTP(w, req)
}
