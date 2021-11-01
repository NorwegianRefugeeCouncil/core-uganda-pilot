package admin

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core/pkg/server/admin/handlers/authn"
	"github.com/nrc-no/core/pkg/server/admin/handlers/clients"
	"github.com/nrc-no/core/pkg/server/admin/handlers/identityprovider"
	"github.com/nrc-no/core/pkg/server/admin/handlers/organization"
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	"golang.org/x/oauth2"
)

type Server struct {
	*generic.Server
	options Options
}

type Options struct {
	options.ServerOptions
	StoreFactory store.Factory
}

func NewServer(options Options) (*Server, error) {

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "admin")
	if err != nil {
		return nil, err
	}

	container := genericServer.Container

	clientsHandler, err := clients.NewHandler()
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

	s := &Server{
		Server:  genericServer,
		options: options,
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) {

	var provider *oidc.Provider
	err := backoff.Retry(func() error {
		var err error
		provider, err = oidc.NewProvider(ctx, s.options.Oidc.Issuer)
		if err != nil {
			s.Logger().WithError(err).Warnf("failed to get oidc provider")
			return err
		}
		return err
	}, backoff.NewExponentialBackOff())

	if err != nil {
		panic(err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     s.options.Oidc.ClientID,
		ClientSecret: s.options.Oidc.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  fmt.Sprintf("%s/oidc/callback", s.options.URLs.Self),
		Scopes:       s.options.Oidc.Scopes,
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID:             s.options.Oidc.ClientID,
		SupportedSigningAlgs: []string{oidc.RS256},
		SkipClientIDCheck:    false,
		SkipExpiryCheck:      false,
		SkipIssuerCheck:      false,
		Now:                  nil,
	})

	authnHandler := authn.NewHandler(
		s.Server.SessionStore(),
		&oauth2Config,
		verifier,
	)

	s.Container.Filter(authn.RestfulAuthnMiddleware(s.SessionStore(), verifier, s.options.URLs.Self, "/"))

	s.Container.Add(authnHandler.WebService())

	s.Server.Start(ctx)
}
