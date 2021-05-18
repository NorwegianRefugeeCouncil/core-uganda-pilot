package options

import (
	"github.com/nrc-no/core/apps/api/pkg/server"
	"github.com/nrc-no/core/apps/api/pkg/server/options"
	"github.com/nrc-no/core/apps/api/pkg/storage/backend"
)

type ServerRunOptions struct {
	Mongo *options.MongoOptions
}

func NewServerRunOptions() *ServerRunOptions {
	s := ServerRunOptions{
		Mongo: options.NewMongoOptions(backend.NewDefaultConfig(server.DefaultMongoPathPrefix, nil)),
	}
	return &s
}
