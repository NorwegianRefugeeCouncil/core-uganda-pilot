package login

import (
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/login/handlers/login"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
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
	sessionStore := genericServer.SessionStore()

	organizationStore := store.NewOrganizationStore(options.StoreFactory)
	idpStore := store.NewIdentityProviderStore(options.StoreFactory)
	loginStore := loginstore.NewStore(options.StoreFactory)

	loginHandler, err := login.NewHandler(
		sessionStore,
		organizationStore,
		idpStore,
		loginStore,
		options.URLs.Self)
	if err != nil {
		return nil, err
	}
	container.Add(loginHandler.WebService())

	return &Server{
		Server: genericServer,
	}, nil
}
