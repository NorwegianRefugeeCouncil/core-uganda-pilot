package auth

import (
	"context"
	"net/http"
)

// SetAuthorizationHeader is a utility method that will add the Access Token to the Authorization header
func SetAuthorizationHeader(ctx context.Context, req *http.Request) {
	if accessToken, ok := ctx.Value(AccessTokenKey).(string); ok {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
}
