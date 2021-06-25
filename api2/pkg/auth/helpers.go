package auth

import (
	"net/http"
)

func GetAuthenticationContext(req *http.Request) AuthenticationContext {
	ctx := req.Context()

	authCtxIntf := ctx.Value(AuthenticationContextKey)
	if authCtxIntf == nil {
		return AuthenticationContext{}
	}
	authCtx, ok := authCtxIntf.(AuthenticationContext)
	if !ok {
		return AuthenticationContext{}
	}

	return authCtx

}

func IsAuthenticatedRequest(req *http.Request) bool {
	return GetAuthenticationContext(req).IsAuthenticated
}
