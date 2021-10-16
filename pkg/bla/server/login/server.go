package login

import (
	"github.com/nrc-no/core/pkg/bla/options"
	"github.com/nrc-no/core/pkg/bla/server/generic"
	"github.com/nrc-no/core/pkg/bla/store"
)

type Server struct {
	*generic.Server
}

type Options struct {
	options.ServerOptions
	StoreFactory store.Factory
}

func NewServer(options Options) (*Server, error) {

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "login")
	if err != nil {
		return nil, err
	}

	return &Server{
		Server: genericServer,
	}, nil
}
