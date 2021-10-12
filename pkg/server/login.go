package server

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/login"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
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

	if err := createOauthClient(ctx, c.HydraAdminClient.Admin, c.HydraTLSClient, cli); err != nil {
		logrus.WithError(err).Errorf("failed to create OAuth2 client")
		return nil, err
	}

	clientCredsCfg := clientcredentials.Config{
		ClientID:     c.LoginClientID,
		ClientSecret: c.LoginClientSecret,
		TokenURL:     c.OAuthTokenEndpoint,
	}

	adminCtx := ctx
	if !c.TLSDisable {
		tlsCli, err := tlsClient(c.TLSCertPath)
		if err != nil {
			return nil, err
		}
		adminCtx = context.WithValue(ctx, oauth2.HTTPClient, tlsCli)
	}
	adminCli := clientCredsCfg.Client(adminCtx)

	loginOptions := &login.ServerOptions{
		GenericServerOptions: genericOptions,
		BCryptCost:           15,
		AdminHTTPClient:      adminCli,
		IAMHost:              c.LoginIAMHost,
		IAMScheme:            c.LoginIAMScheme,
		TemplateDirectory:    c.LoginTemplateDirectory,
	}

	loginServer, err := login.NewServer(ctx, loginOptions)
	if err != nil {
		logrus.WithError(err).Errorf("failed to create login server")
		return nil, err
	}

	return loginServer, nil
}
