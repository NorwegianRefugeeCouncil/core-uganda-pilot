package login

import (
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/login/handlers/login"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
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

	container := genericServer.Container

	organizationStore := store.NewOrganizationStore(options.StoreFactory)
	idpStore := store.NewIdentityProviderStore(options.StoreFactory)
	loginHandler, err := login.NewHandler(organizationStore, idpStore)
	if err != nil {
		return nil, err
	}
	container.Add(loginHandler.WebService())

	return &Server{
		Server: genericServer,
	}, nil
}
