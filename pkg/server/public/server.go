package public

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/generic"
	authn2 "github.com/nrc-no/core/pkg/server/handlers/authn"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/server/public/handlers/database"
	"github.com/nrc-no/core/pkg/server/public/handlers/folder"
	"github.com/nrc-no/core/pkg/server/public/handlers/form"
	"github.com/nrc-no/core/pkg/server/public/handlers/record"
	store2 "github.com/nrc-no/core/pkg/store"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Server struct {
	*generic.Server
	options Options
}

type Options struct {
	options.ServerOptions
	StoreFactory store2.Factory
}

func NewServer(options Options) (*Server, error) {

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "public")
	if err != nil {
		return nil, err
	}

	container := genericServer.Container

	databaseStore := store2.NewDatabaseStore(options.StoreFactory)
	databaseHandler := database.NewHandler(databaseStore)
	container.Add(databaseHandler.WebService())

	folderStore := store2.NewFolderStore(options.StoreFactory)
	folderHandler := folder.NewHandler(folderStore)
	container.Add(folderHandler.WebService())

	formStore := store2.NewFormStore(options.StoreFactory)
	formHandler := form.NewHandler(formStore)
	container.Add(formHandler.WebService())

	recordStore := store2.NewRecordStore(options.StoreFactory, formStore)
	recordHandler := record.NewHandler(recordStore)
	container.Add(recordHandler.WebService())

	s := &Server{
		options: options,
		Server:  genericServer,
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) {

	l := logging.NewLogger(ctx)

	var provider *oidc.Provider
	err := backoff.Retry(func() error {
		var err error
		provider, err = oidc.NewProvider(ctx, s.options.Oidc.Issuer)
		if err != nil {
			l.With(zap.Error(err)).Warn("failed to get oidc provider")
			return err
		}
		return err
	}, backoff.NewExponentialBackOff())

	if err != nil {
		panic(err)
	}

	oauth2Config := &oauth2.Config{
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

	authnHandler := authn2.NewHandler(
		"core-app-session",
		s.options.Oidc.RedirectURI,
		s.Server.SessionStore(),
		oauth2Config,
		verifier,
	)

	s.Container.Filter(authn2.RestfulAuthnMiddleware(
		s.SessionStore(),
		oauth2Config,
		verifier,
		s.options.URLs.Self,
		"core-app-session"),
	)

	s.Container.Add(authnHandler.WebService())

	s.Server.Start(ctx)
}
