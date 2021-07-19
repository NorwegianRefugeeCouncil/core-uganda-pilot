package testutils

import (
	"github.com/nrc-no/core/pkg/rest"
)

func SetXAuthenticatedUserSubject(port string) *rest.RESTConfig {
	return &rest.RESTConfig{
		Scheme: "http",
		Host:   "localhost:" + port,
		Headers: map[string][]string{
			"X-Authenticated-User-Subject": {"mock-auth-user"},
		},
	}
}
