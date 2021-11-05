package authn

import (
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/utils/sets"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

var unprotectedPaths = sets.NewString(
	"/oidc/login",
	"/oidc/renew",
	"/oidc/callback",
	"/oidc/session",
)

func RestfulAuthnMiddleware(
	sessionStore sessions.Store,
	oauth2Config *oauth2.Config,
	verifier *oidc.IDTokenVerifier,
	selfURL string,
	sessionKey string,
	tokenVerifier *oidc.IDTokenVerifier,
) func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

		ctx := logging.WithMiddleware(request.Request.Context(), "auth")
		l := logging.NewLogger(ctx)

		l.Debug("authn callback called")

		if unprotectedPaths.Has(request.Request.URL.Path) {
			l.Debug("ignoring authentication because request is for an unprotected path", zap.String("path", request.Request.URL.Path))
			chain.ProcessFilter(request, response)
			return
		}

		redirectToLogin := func(req *restful.Request, response *restful.Response) {
			l.Debug("redirecting to login")
			response.Header().Set("Location", fmt.Sprintf("%s/oidc/login", selfURL))
			if req.HeaderParameter("Accept") == "application/json" {
				utils.JSONResponse(response.ResponseWriter, http.StatusUnauthorized, meta.NewUnauthorized("not logged in"))
			} else {
				utils.JSONResponse(response.ResponseWriter, http.StatusTemporaryRedirect, meta.NewUnauthorized("not logged in"))
			}
			return
		}

		authorizationHeader := request.HeaderParameter("Authorization")
		if len(authorizationHeader) > 0 {
			l.Debug("authorization header present. introspecting...")
			if authorizationHeader[:7] != "Bearer " {
				l.Error("malformed authorization header")
				redirectToLogin(request, response)
				return
			}
			parts := strings.Split(authorizationHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]

				_, err := tokenVerifier.Verify(ctx, token)
				if err != nil {
					l.Error("token is invalid", zap.Error(err))
					redirectToLogin(request, response)
					return
				}

				l.Debug("valid authorization token found")
				chain.ProcessFilter(request, response)
				return
			} else {
				l.Error("malformed authorization header")
				redirectToLogin(request, response)
				return
			}
		}

		l.Debug("getting user session")

		userSession, err := sessionStore.Get(request.Request, sessionKey)
		if err != nil {
			l.Error("failed to get user session", zap.Error(err))
			redirectToLogin(request, response)
			return
		}

		l.Debug("getting token from session")

		previousTokenFromSession, ok := tokenFromSession(userSession)
		if !ok {
			l.Warn("redirecting to login as previous token was not found in session")
			redirectToLogin(request, response)
			return
		}

		previousRefreshToken := previousTokenFromSession.RefreshToken
		previousAccessToken := previousTokenFromSession.AccessToken
		tokenSource := oauth2Config.TokenSource(ctx, previousTokenFromSession.Token)

		l.Debug("refreshing token if necessary")

		newToken, err := tokenSource.Token()
		if err != nil {
			l.Warn("failed to get token. perhaps token was expired", zap.Error(err))
			redirectToLogin(request, response)
			return
		}

		l.Debug("checking if token was refreshed")

		if newToken.RefreshToken == previousRefreshToken && newToken.AccessToken == previousAccessToken {
			l.Debug("refresh and access token are identical. proceeding")
			chain.ProcessFilter(request, response)
			return
		}

		l.Debug("verifying token")

		rawIdToken, idToken, err := verifyToken(ctx, verifier, newToken)
		if err != nil {
			l.Error("failed to verify new token", zap.Error(err))
			redirectToLogin(request, response)
			return
		}

		l.Debug("populating session with new token")

		// here we re-store the token because it was refreshed
		userSession.Values[constants.SessionIDToken] = rawIdToken
		userSession.Values[constants.SessionAccessToken] = newToken.AccessToken
		userSession.Values[constants.SessionRefreshToken] = newToken.RefreshToken
		userSession.Values[constants.SessionTokenExpiry] = newToken.Expiry
		var userProfile Claims
		if err := idToken.Claims(&userProfile); err != nil {
			l.Error("failed to unmarshal token claims", zap.Error(err))
			redirectToLogin(request, response)
			return
		}

		l.Debug("unmarshaled token claims")

		l.Debug("saving user session")

		userSession.Values[constants.SessionProfile] = userProfile
		if err := userSession.Save(request.Request, response.ResponseWriter); err != nil {
			l.Error("failed to save user session with new token", zap.Error(err))
			redirectToLogin(request, response)
			return
		}

		l.Debug("done")

		chain.ProcessFilter(request, response)
		return
	}
}

type StoredToken struct {
	*oauth2.Token
	IDToken string
}
