package authorization

import (
	"github.com/nrc-no/core/pkg/api/meta"
	"net/http"
	"strings"
)

func ExtractBearerToken(req *http.Request) (string, error) {

	if req.Header == nil {
		return "", meta.NewUnauthorized("missing header")
	}

	authHeaderParts := req.Header["Authorization"]

	if len(authHeaderParts) == 0 {
		return "", meta.NewUnauthorized("missing token")
	}

	if len(authHeaderParts) != 1 {
		return "", meta.NewUnauthorized("invalid authorization header")
	}

	authHeader := authHeaderParts[0]

	if len(authHeader) == 0 {
		return "", meta.NewUnauthorized("no Authorization header in request")
	}
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" || len(parts[1]) == 0 {
		return "", meta.NewUnauthorized("malformed Authorization header")
	}

	return parts[1], nil
}
