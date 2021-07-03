package server

import (
	"context"
	"github.com/nrc-no/core/pkg/apps/login"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/ory/hydra-client-go/models"
	"golang.org/x/oauth2/clientcredentials"
)

func (c CompletedOptions) CreateLoginServer(ctx context.Context, genericOptions *server.GenericServerOptions) (*login.Server, error) {
	cli := &models.OAuth2Client{
		ClientID:                c.LoginClientID,
		ClientName:              c.LoginClientName,
		ClientSecret:            c.LoginClientSecret,
		GrantTypes:              []string{"client_credentials"},
		ResponseTypes:           []string{"token", "refresh_token"},
		TokenEndpointAuthMethod: "client_secret_post",
	}

	if err := createOauthClient(ctx, c.HydraAdminClient.Admin, cli); err != nil {
		return nil, err
	}

	clientCredsCfg := clientcredentials.Config{
		ClientID:     c.LoginClientID,
		ClientSecret: c.LoginClientSecret,
		TokenURL:     c.OAuthTokenEndpoint,
	}
	adminCli := clientCredsCfg.Client(ctx)

	loginOptions := &login.ServerOptions{
		GenericServerOptions: genericOptions,
		BCryptCost:           15,
		AdminHTTPClient:      adminCli,
		IAMHost:              c.LoginIAMHost,
		IAMScheme:            c.LoginIAMScheme,
	}

	loginServer, err := login.NewServer(ctx, loginOptions)
	if err != nil {
		return nil, err
	}

	return loginServer, nil
}
