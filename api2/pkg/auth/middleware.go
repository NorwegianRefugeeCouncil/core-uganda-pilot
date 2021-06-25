package auth

import (
	"context"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strings"
	"time"
)

type Roles struct {
	Roles []string `json:"roles"`
}

type ResourceAccessMap map[string]Roles

type Claims struct {
	jwt.Claims        `json:",inline"`
	PreferredUsername string            `json:"preferred_username"`
	EmailVerified     bool              `json:"email_verified"`
	Scope             string            `json:"scope"`
	FamilyName        string            `json:"family_name"`
	GivenName         string            `json:"given_name"`
	Name              string            `json:"name"`
	RealmAccess       Roles             `json:"realm_access"`
	ResourceAccess    ResourceAccessMap `json:"resource_access"`
}

type UserInfo struct {
	Subject           string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}

type AuthenticationContext struct {
	Claims          Claims
	IsAuthenticated bool
	AccessToken     string
}

const AuthenticationContextKey = "authentication_context"

func (h *Handler) Authenticate() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			ctx := req.Context()

			// Do not redirect when the user is trying to authenticate
			shouldRedirect := !strings.HasPrefix(req.URL.Path, "/auth/")

			doNext := func(req *http.Request, isAuthenticated bool) {
				if !isAuthenticated {
					ctx = req.Context()
					ctx = context.WithValue(ctx, AuthenticationContextKey, AuthenticationContext{
						IsAuthenticated: false,
					})
					req = req.WithContext(ctx)
				}

				if !isAuthenticated && shouldRedirect {
					http.Redirect(w, req, "/auth/login", http.StatusTemporaryRedirect)
					return
				} else {
					next.ServeHTTP(w, req)
					return
				}
			}

			// Retrieve from store if available
			accessToken := h.Store.GetString(ctx, AccessTokenKey)

			// If not, fallback to Authorization header
			if len(accessToken) == 0 {
				bearerStr := req.Header.Get("Authorization")
				if strings.HasPrefix(bearerStr, "Bearer ") {
					accessToken = strings.TrimPrefix(bearerStr, "Bearer ")
				}
			}

			// If we don't have an access token via the session or via
			// the Authorization header, then assume we are not authenticated
			if len(accessToken) == 0 {
				doNext(req, false)
				return
			}

			// Retrieve Keycloak public keys for token verification
			// TODO: cache this...
			jwkSet, err := h.KeycloakClient.GetPublicKeys()
			if err != nil {
				doNext(req, false)
				return
			}

			// That should not happen, but whatever
			if jwkSet == nil {
				doNext(req, false)
				return
			}

			// Make sure that we have some keys in the jwkSet
			// TODO: figure out what happens when we have multiple keys
			if len(jwkSet.Keys) == 0 {
				doNext(req, false)
				return
			}

			// Parse the token
			jwtToken, err := jwt.ParseSigned(accessToken)
			if err != nil {
				doNext(req, false)
				return
			}

			if jwtToken == nil {
				doNext(req, false)
				return
			}

			jwkKey := jwkSet.Keys[0].Key

			// Unmarshal the access token claims
			var claims Claims
			if err := jwtToken.Claims(jwkKey, &claims); err != nil {
				doNext(req, false)
				return
			}

			// Validate the token
			if err := claims.Validate(jwt.Expected{
				Issuer:   "http://localhost:8080/auth/realms/nrc",
				Audience: jwt.Audience{"account"},
				Time:     time.Now(),
			}); err != nil {
				doNext(req, false)
				return
			}

			// Check if the party corresponding to the user exists

			// Store the AuthenticationContext in the context,
			// allowing downstream usage
			authCtx := AuthenticationContext{
				Claims:          claims,
				IsAuthenticated: true,
				AccessToken:     accessToken,
			}
			ctx = context.WithValue(ctx, AuthenticationContextKey, authCtx)
			req = req.WithContext(ctx)

			// At this point, assume we are authenticated
			doNext(req, true)
			return

		})
	}
}
