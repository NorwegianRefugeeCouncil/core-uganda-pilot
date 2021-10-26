package authn

import (
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
	"strings"
)

func RestfulAuthnMiddleware(
	sessionStore sessions.Store,
	verifier *oidc.IDTokenVerifier,
	selfURL string,
	defaultRedirectURL string,
) func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {

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

			if userSession, err := sessionStore.Get(req.Request, constants.SessionKey); err == nil {
				requestedURL := req.Request.URL.String()
				if !strings.HasPrefix(requestedURL, loginPath) && !strings.HasPrefix(requestedURL, callbackPath) {
					userSession.Values[constants.SessionDesiredURL] = requestedURL
					_ = userSession.Save(req.Request, response.ResponseWriter)
				} else {
					userSession.Values[constants.SessionDesiredURL] = defaultRedirectURL
				}
			}

			response.Header().Set("Location", fmt.Sprintf("%s/oidc/login", selfURL))
			response.WriteHeader(http.StatusTemporaryRedirect)
			utils.JSONResponse(response.ResponseWriter, http.StatusTemporaryRedirect, meta.NewUnauthorized("not logged in"))
			return
		}

		tokenFromHeader, hasTokenFromHeader := idTokenFromHeader(request.Request)
		tokenFromSession, hasTokenFromSession := idTokenFromSession(request.Request, sessionStore)

		var rawIDToken string
		if hasTokenFromHeader {
			rawIDToken = tokenFromHeader
		} else if hasTokenFromSession {
			rawIDToken = tokenFromSession
		} else {
			redirectToLogin(request, response)
			return
		}

		_, err := verifier.Verify(request.Request.Context(), rawIDToken)
		if err != nil {
			redirectToLogin(request, response)
			return
		}

		chain.ProcessFilter(request, response)
	}
}

func idTokenFromHeader(req *http.Request) (string, bool) {

	header := req.Header.Get("Authorization")
	if len(header) == 0 {
		return "", false
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", false
	}

	if parts[0] != "Bearer " {
		return "", false
	}

	if len(parts[1]) == 0 {
		return "", false
	}

	return parts[1], true

}

func idTokenFromSession(req *http.Request, sessionStore sessions.Store) (string, bool) {
	userSession, err := sessionStore.Get(req, constants.SessionKey)
	if err != nil {
		return "", false
	}

	idTokenIntf, ok := userSession.Values[constants.SessionIDToken]
	if !ok {
		return "", false
	}

	rawIDToken, ok := idTokenIntf.(string)
	if !ok {
		return "", false
	}

	return rawIDToken, true
}
