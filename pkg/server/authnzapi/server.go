package authnzapi

import (
	"context"
	"github.com/nrc-no/core/pkg/server/authnzapi/handlers/clients"
	"github.com/nrc-no/core/pkg/server/authnzapi/handlers/identity"
	"github.com/nrc-no/core/pkg/server/authnzapi/handlers/identityprovider"
	"github.com/nrc-no/core/pkg/server/authnzapi/handlers/organization"
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/client/admin"
)

type Server struct {
	*generic.Server
	options Options
}

type Options struct {
	options.ServerOptions
	StoreFactory store.Factory
	HydraAdmin   admin.ClientService
}

func NewServer(options Options) (*Server, error) {

	hydraAdmin := options.HydraAdmin

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "authnz-api")
	if err != nil {
		return nil, err
	}

	container := genericServer.GoRestfulContainer

	clientsHandler, err := clients.NewHandler(hydraAdmin)
	if err != nil {
		return nil, err
	}
	container.Add(clientsHandler.WebService())

	organizationStore := store.NewOrganizationStore(options.StoreFactory)
	// idpStore := store.NewIdentityProviderStore(options.StoreFactory)
	organizationsHandler := organization.NewHandler(organizationStore)
	if err != nil {
		return nil, err
	}
	container.Add(organizationsHandler.WebService())

	idpStore := store.NewIdentityProviderStore(options.StoreFactory)
	idpHandler := identityprovider.NewHandler(idpStore)
	container.Add(idpHandler.WebService())

	identityStore := store.NewIdentityStore(options.StoreFactory)
	identityHandler := identity.NewHandler(identityStore)
	container.Add(identityHandler.WebService())

	s := &Server{
		Server:  genericServer,
		options: options,
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) {
	s.Server.Start(ctx)
}
