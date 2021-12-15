package authnzbouncer

import (
	"context"
	"github.com/nrc-no/core/pkg/server/authnzbouncer/authenticators"
	"github.com/nrc-no/core/pkg/server/authnzbouncer/authorizers"
	"github.com/nrc-no/core/pkg/server/authnzbouncer/handlers"
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/client/public"
)

type Server struct {
	*generic.Server
	options Options
}

type Options struct {
	options.ServerOptions
	HydraAdmin  admin.ClientService
	HydraPublic public.ClientService
}

func NewServer(options Options) (*Server, error) {
	genericServer, err := generic.NewGenericServer(options.ServerOptions, "authnz-bouncer")
	if err != nil {
		return nil, err
	}
	authenticator := authenticators.NewHydraAuthenticator(options.HydraPublic)
	authorizer := authorizers.NewHydraAuthorizer(options.HydraAdmin)
	genericServer.NonGoRestfulMux.PathPrefix("/").Handler(handlers.HandleAuth(authenticator, authorizer))
	s := &Server{
		options: options,
		Server:  genericServer,
	}
	return s, nil
}

func (s *Server) Start(ctx context.Context) {
	s.Server.Start(ctx)
}
