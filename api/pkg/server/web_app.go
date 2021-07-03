package server

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core/pkg/apps/webapp"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/ory/hydra-client-go/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func (c CompletedOptions) CreateWebAppServer(ctx context.Context, genericOptions *server.GenericServerOptions) (*webapp.Server, error) {

	cli := &models.OAuth2Client{
		ClientID:     c.WebAppClientID,
		ClientName:   c.WebAppClientName,
		ClientSecret: c.WebAppClientSecret,
		GrantTypes: []string{
			"client_credentials",
			"authorization_code",
			"id_token",
			"access_token",
			"refresh_token",
		},
		RedirectUris: []string{
			c.BaseURL + "/callback",
		},
		ResponseTypes: []string{
			"token",
			"code",
		},
		Scope:                   "openid profile",
		TokenEndpointAuthMethod: "client_secret_post",
		PostLogoutRedirectUris: []string{
			c.BaseURL,
		},
	}

	if err := createOauthClient(ctx, c.HydraAdminClient.Admin, cli); err != nil {
		return nil, err
	}

	clientCredsCfg := clientcredentials.Config{
		ClientID:     c.WebAppClientID,
		ClientSecret: c.WebAppClientSecret,
		TokenURL:     c.OAuthTokenEndpoint,
	}
	adminCli := clientCredsCfg.Client(ctx)

	oidcVerifier := c.OIDCProvider.Verifier(&oidc.Config{
		ClientID: c.WebAppClientID,
	})

	oauth2Config := &oauth2.Config{
		ClientID:     c.WebAppClientID,
		ClientSecret: c.WebAppClientSecret,
		Endpoint:     c.OIDCProvider.Endpoint(),
		RedirectURL:  c.BaseURL + "/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	webAppOptions := &webapp.ServerOptions{
		GenericServerOptions: genericOptions,
		TemplateDirectory:    c.WebAppTemplateDirectory,
		BaseURL:              c.BaseURL + c.WebAppBasePath,
		IAMHost:              c.WebAppIAMHost,
		IAMScheme:            c.WebAppIAMScheme,
		CMSHost:              c.WebAppCMSHost,
		CMSScheme:            c.WebAppCMSScheme,
		AdminHTTPClient:      adminCli,
		IDTokenVerifier:      oidcVerifier,
		OAuth2Config:         oauth2Config,
	}

	webappServer, err := webapp.NewServer(webAppOptions)
	if err != nil {
		panic(err)
	}

	return webappServer, nil

}
