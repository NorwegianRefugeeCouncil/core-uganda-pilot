package auth

import (
	"net/http"
)

func IsAuthenticatedRequest(req *http.Request) bool {
	ctx := req.Context()
	isLoggedIn := ctx.Value(IsLoggedInKey)
	if isLoggedIn == nil {
		return false
	}
	isLoggedInBool, ok := isLoggedIn.(bool)
	if !ok {
		return false
	}
	return isLoggedInBool
}
