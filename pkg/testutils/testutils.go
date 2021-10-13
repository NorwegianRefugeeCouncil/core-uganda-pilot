package testutils

import (
	"github.com/nrc-no/core/pkg/rest"
	uuid "github.com/satori/go.uuid"
)

func SetXAuthenticatedUserSubject(port string) *rest.Config {
	return &rest.Config{
		Scheme: "http",
		Host:   "localhost:" + port,
		Headers: map[string][]string{
			"X-Authenticated-User-Subject": {uuid.NewV4().String()},
		},
	}
}