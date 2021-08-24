package server

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core/pkg/apps/webapp"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

func (c CompletedOptions) CreateWebAppServer(ctx context.Context, genericOptions *server.GenericServerOptions) (*webapp.Server, error) {

	l := logrus.WithField("server", "webapp")
	l.Infof("creating webapp server")

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

	l.Infof("creating webapp oauth2 client")
	if err := createOauthClient(ctx, c.HydraAdminClient.Admin, c.HydraTLSClient, cli); err != nil {
		l.WithError(err).Errorf("failed to create webapp oauth2 client")
		return nil, err
	}

	l.Infof("configuring HTTP client for internal traffic")
	httpClient := http.DefaultClient
	if !c.TLSDisable {
		var err error
		httpClient, err = tlsClient(c.TLSCertPath)
		if err != nil {
			l.WithError(err).Errorf("failed to create internal HTTP client")
			return nil, err
		}
	}

	l.Infof("creating administrative HTTP client for privileged traffic")
	clientCredsCfg := clientcredentials.Config{
		ClientID:     c.WebAppClientID,
		ClientSecret: c.WebAppClientSecret,
		TokenURL:     c.OAuthTokenEndpoint,
	}
	adminCli := clientCredsCfg.Client(context.WithValue(ctx, oauth2.HTTPClient, httpClient))

	l.Infof("creating private oauth2 configuration")
	privateOauth2Config := &oauth2.Config{
		ClientID:     c.WebAppClientID,
		ClientSecret: c.WebAppClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.HydraPublicURL + "/oauth2/auth",
			TokenURL: c.HydraPublicURL + "/oauth2/token",
		},
		RedirectURL: c.BaseURL + "/callback",
		Scopes:      []string{oidc.ScopeOpenID, "profile"},
	}

	l.Infof("creating public oauth2 configuration")
	publicOauth2Config := &oauth2.Config{
		ClientID:     c.WebAppClientID,
		ClientSecret: c.WebAppClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  *c.OpenIdConf.Payload.AuthorizationEndpoint,
			TokenURL: *c.OpenIdConf.Payload.TokenEndpoint,
		},
		RedirectURL: c.BaseURL + "/callback",
		Scopes:      []string{oidc.ScopeOpenID, "profile"},
	}

	l.Infof("creating oidc verifier")
	jwks := oidc.NewRemoteKeySet(ctx, c.OAuthJwksURI)
	oidcVerifier := oidc.NewVerifier(*c.OpenIdConf.Payload.Issuer, jwks, &oidc.Config{
		ClientID:             c.WebAppClientID,
		SupportedSigningAlgs: c.OAuthIDTokenSigningAlgs,
	})

	webAppOptions := &webapp.ServerOptions{
		GenericServerOptions: genericOptions,
		TemplateDirectory:    c.WebAppTemplateDirectory,
		StaticDic:            c.WebAppStaticDir,
		BaseURL:              c.BaseURL + c.WebAppBasePath,
		IAMHost:              c.WebAppIAMHost,
		IAMScheme:            c.WebAppIAMScheme,
		CMSHost:              c.WebAppCMSHost,
		CMSScheme:            c.WebAppCMSScheme,
		AdminHTTPClient:      adminCli,
		IDTokenVerifier:      oidcVerifier,
		PrivateOAuth2Config:  privateOauth2Config,
		PublicOauth2Config:   publicOauth2Config,
		HydraHTTPClient:      httpClient,
		IAMHTTPClient:        httpClient,
		CMSHTTPClient:        httpClient,
		LoginHTTPClient:      httpClient,
	}

	webappServer, err := webapp.NewServer(webAppOptions)
	if err != nil {
		return nil, err
	}

	return webappServer, nil

}
