package authn

import (
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/utils"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func RestfulAuthnMiddleware(
	sessionStore sessions.Store,
	verifier *oidc.IDTokenVerifier,
	oauth2Config *oauth2.Config,
	selfURL string,
	sessionKey string,
) func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

		ctx := request.Request.Context()

		loginPath := "/oidc/login"
		callbackPath := "/oidc/callback"

		if request.Request.URL.Path == loginPath {
			chain.ProcessFilter(request, response)
			return
		}

		if request.Request.URL.Path == callbackPath {
			chain.ProcessFilter(request, response)
			return
		}

		redirectToLogin := func(req *restful.Request, response *restful.Response) {
			response.Header().Set("Location", fmt.Sprintf("%s/oidc/login", selfURL))
			if req.HeaderParameter("Accept") == "application/json" {
				utils.JSONResponse(response.ResponseWriter, http.StatusUnauthorized, meta.NewUnauthorized("not logged in"))
			} else {
				utils.JSONResponse(response.ResponseWriter, http.StatusTemporaryRedirect, meta.NewUnauthorized("not logged in"))
			}
			return
		}

		if userSession, err := sessionStore.Get(request.Request, sessionKey); err == nil {
			if previousTokenFromSession, ok := tokenFromSession(userSession); ok {
				tokenSource := oauth2Config.TokenSource(ctx, previousTokenFromSession)
				if newToken, err := tokenSource.Token(); err == nil {
					// here we re-store the token in case oauth2 refreshed it.
					userSession.Values[constants.SessionAccessToken] = newToken.AccessToken
					userSession.Values[constants.SessionRefreshToken] = newToken.RefreshToken
					userSession.Values[constants.SessionTokenExpiry] = newToken.Expiry
					if err := userSession.Save(request.Request, response.ResponseWriter); err == nil {
						chain.ProcessFilter(request, response)
						return
					}
				}
			}
		}

		redirectToLogin(request, response)
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

func tokenFromSession(userSession *sessions.Session) (*oauth2.Token, bool) {
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
	return &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		RefreshToken: refreshToken,
		Expiry:       *tokenExpiry,
	}, true
}
