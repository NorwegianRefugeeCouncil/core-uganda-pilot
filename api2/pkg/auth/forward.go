package auth

import (
	"context"
	"net/http"
)

// SetAuthorizationHeader is a utility method that will add the Access Token to the Authorization header
func SetAuthorizationHeader(ctx context.Context, req *http.Request) {
	authCtx := GetAuthenticationContext(req)
	if len(authCtx.AccessToken) > 0 {
		req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
	}
}
