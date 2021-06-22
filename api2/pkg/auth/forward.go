package auth

import (
	"context"
	"net/http"
)

func Forward(ctx context.Context, req *http.Request) {
	if accessToken, ok := ctx.Value(AccessTokenKey).(string); ok {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
}
