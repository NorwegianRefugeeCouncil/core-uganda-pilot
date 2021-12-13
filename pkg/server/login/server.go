package login

import (
	"context"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/login/handlers/login"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/client/admin"
	"go.uber.org/zap"
)

type Server struct {
	*generic.Server
}

type Options struct {
	options.ServerOptions
	StoreFactory store.Factory
	HydraAdmin   admin.ClientService
}

func NewServer(ctx context.Context, options Options) (*Server, error) {

	l := logging.NewLogger(ctx)

	l.Info("using configuration",
		zap.String("host", options.Host),
		zap.Int("port", options.Port),
		zap.String("urls.self", options.URLs.Self),
	)

	hydraAdmin := options.HydraAdmin

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "login")
	if err != nil {
		return nil, err
	}

	container := genericServer.GoRestfulContainer
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

	for _, service := range loginHandler.WebServices() {
		container.Add(service)
	}

	return &Server{
		Server: genericServer,
	}, nil
}
