package auth

import (
	"context"
	jwt2 "gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) Authenticate() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			ctx := req.Context()

			shouldRedirect := !strings.HasPrefix(req.URL.Path, "/auth/")
			doNext := func(req *http.Request, isAuthenticated bool) {

				ctx = context.WithValue(ctx, IsLoggedInKey, isAuthenticated)
				req = req.WithContext(ctx)

				if !isAuthenticated && shouldRedirect {
					http.Redirect(w, req, "/auth/login", http.StatusTemporaryRedirect)
				}
				next.ServeHTTP(w, req)
			}

			// Retrieve from store if available
			accessToken := h.Store.GetString(ctx, AccessTokenKey)
			isInStore := len(accessToken) != 0

			// If not, fallback to Authorization header
			if len(accessToken) == 0 {
				bearerStr := req.Header.Get("Authorization")
				if strings.HasPrefix(bearerStr, "Bearer ") {
					accessToken = strings.TrimPrefix(bearerStr, "Bearer ")
				}
			}

			if len(accessToken) == 0 {
				doNext(req, false)
				return
			}

			jkws, err := h.KeycloakClient.GetPublicKeys()
			if err != nil {
				doNext(req, false)
				return
			}

			jwtToken, err := jwt2.ParseSigned(accessToken)
			if err != nil {
				doNext(req, false)
				return
			}

			type claims struct {
				jwt2.Claims       `json:",inline"`
				PreferredUsername string `json:"preferred_username"`
				EmailVerified     bool   `json:"email_verified"`
			}

			var a claims

			if err := jwtToken.Claims(jkws.Keys[0].Key, &a); err != nil {
				doNext(req, false)
				return
			}

			if err := a.Validate(jwt2.Expected{
				Issuer:   "http://localhost:8080/auth/realms/nrc",
				Audience: jwt2.Audience{"account"},
				Time:     time.Now(),
			}); err != nil {
				doNext(req, false)
				return
			}

			if !isInStore {
				h.Store.Put(ctx, AccessTokenKey, accessToken)
			}

			ctx = context.WithValue(ctx, AccessTokenKey, accessToken)
			ctx = context.WithValue(ctx, IsLoggedInKey, true)
			req = req.WithContext(ctx)

			doNext(req, true)
			return

		})
	}
}
