package login

import (
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/login/handlers/login"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/client/admin"
)

type Server struct {
	*generic.Server
}

type Options struct {
	options.ServerOptions
	StoreFactory store.Factory
	HydraAdmin   admin.ClientService
}

func NewServer(options Options) (*Server, error) {

	hydraAdmin := options.HydraAdmin

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "login")
	if err != nil {
		return nil, err
	}

	container := genericServer.Container
	sessionStore := genericServer.SessionStore()

	organizationStore := store.NewOrganizationStore(options.StoreFactory)
	idpStore := store.NewIdentityProviderStore(options.StoreFactory)
	loginStore := loginstore.NewStore(options.StoreFactory)

	loginHandler, err := login.NewHandler(
		sessionStore,
		organizationStore,
		idpStore,
		loginStore,
		options.URLs.Self,
		hydraAdmin)
	if err != nil {
		return nil, err
	}
	container.Add(loginHandler.WebService())

	return &Server{
		Server: genericServer,
	}, nil
}
