package auth

import (
	"context"
	jwt2 "gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strings"
)

func (h *Handler) Authenticate() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			ctx := req.Context()

			shouldRedirect := !strings.HasPrefix(req.URL.Path, "/auth/")
			doNext := func(req *http.Request, isAuthenticated bool) {
				if !isAuthenticated && shouldRedirect {
					http.Redirect(w, req, "/auth/login", http.StatusTemporaryRedirect)
				}
				next.ServeHTTP(w, req)
			}

			// Retrieve from store if available
			accessToken := h.Store.GetString(ctx, AccessTokenKey)

			// If not, fallback to Authorization header
			if len(accessToken) == 0 {

				bearerStr := req.Header.Get("Authorization")
				if strings.HasPrefix(bearerStr, "Bearer ") {
					accessToken = strings.TrimPrefix(bearerStr, "Bearer ")
				} else {
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

				var a = struct {
				}{}

				if err := jwtToken.Claims(jkws.Keys[0].Key, &a); err != nil {
					doNext(req, false)
					return
				}

				h.Store.Put(ctx, AccessTokenKey, accessToken)

			}

			if len(accessToken) == 0 {
				doNext(req, false)
				return
			}

			ctx = context.WithValue(ctx, AccessTokenKey, accessToken)
			req = req.WithContext(ctx)

			doNext(req, true)
			return

		})
	}
}
