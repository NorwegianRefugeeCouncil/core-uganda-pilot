package authorization

import (
	"github.com/nrc-no/core/pkg/api/meta"
	"net/http"
	"strings"
)

func ExtractBearerToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")

	if len(authHeader) == 0 {
		return "", meta.NewUnauthorized("no Authorization header in request")
	}
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" || len(parts[1]) == 0 {
		return "", meta.NewUnauthorized("malformed Authorization header")
	}

	return parts[1], nil
}
