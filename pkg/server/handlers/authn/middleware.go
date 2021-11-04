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
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func RestfulAuthnMiddleware(
	sessionStore sessions.Store,
	oauth2Config *oauth2.Config,
	verifier *oidc.IDTokenVerifier,
	selfURL string,
	sessionKey string,
) func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

		ctx := logging.WithMiddleware(request.Request.Context(), "auth")
		l := logging.NewLogger(ctx)

		loginPath := "/oidc/login"
		callbackPath := "/oidc/callback"

		l.Debug("authn callback called")

		if request.Request.URL.Path == loginPath {
			l.Debug("ignoring authentication because request is for the /oidc/login path")
			chain.ProcessFilter(request, response)
			return
		}

		if request.Request.URL.Path == callbackPath {
			l.Debug("ignoring authentication because request is for the /oidc/callback path")
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

func strFromSession(session *sessions.Session, key string) (string, bool) {
	valueIntf, ok := session.Values[key]
	if !ok {
		return "", false
	}
	value, ok := valueIntf.(string)
	if !ok {
		return "", false
	}
	return value, true
}

func timeFromSession(session *sessions.Session, key string) (*time.Time, bool) {
	valueIntf, ok := session.Values[key]
	if !ok {
		return nil, false
	}
	value, ok := valueIntf.(*time.Time)
	if !ok {
		return nil, false
	}
	return value, true
}

type StoredToken struct {
	*oauth2.Token
	IDToken string
}

func tokenFromSession(userSession *sessions.Session) (*StoredToken, bool) {
	idToken, ok := strFromSession(userSession, constants.SessionIDToken)
	if !ok {
		return nil, false
	}
	accessToken, ok := strFromSession(userSession, constants.SessionAccessToken)
	if !ok {
		return nil, false
	}
	refreshToken, ok := strFromSession(userSession, constants.SessionRefreshToken)
	if !ok {
		return nil, false
	}
	tokenType, ok := strFromSession(userSession, constants.SessionTokenType)
	if !ok {
		return nil, false
	}
	tokenExpiry, ok := timeFromSession(userSession, constants.SessionTokenExpiry)
	if !ok {
		return nil, false
	}
	return &StoredToken{
		Token: &oauth2.Token{
			AccessToken:  accessToken,
			TokenType:    tokenType,
			RefreshToken: refreshToken,
			Expiry:       *tokenExpiry,
		},
		IDToken: idToken,
	}, true
}
