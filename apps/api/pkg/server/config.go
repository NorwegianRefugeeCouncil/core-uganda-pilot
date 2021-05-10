package server

import "github.com/nrc-no/core/apps/api/pkg/registry/generic"

type Config struct {
	RESTOptionsGetter generic.RESTOptionsGetter
}
