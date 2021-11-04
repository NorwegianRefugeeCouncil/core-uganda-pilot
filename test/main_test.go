package test

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nrc-no/core/pkg/client"
	"github.com/nrc-no/core/pkg/rest"
	hydra "github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"testing"
)

type Suite struct {
	suite.Suite
	oauth2Cfg clientcredentials.Config
}

func (s *Suite) Client(ctx context.Context) *http.Client {
	return s.oauth2Cfg.Client(ctx)
}

func (s *Suite) PublicClient(ctx context.Context) client.Client {
	return client.NewClientFromConfig(rest.Config{
		Scheme:     "http",
		Host:       "localhost:9000",
		HTTPClient: s.Client(ctx),
	})
}

func (s *Suite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hydraAdmin := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: "localhost:4445"}).Admin
	_, err := hydraAdmin.DeleteOAuth2Client(&admin.DeleteOAuth2ClientParams{Context: ctx, ID: "integration-testing"})
	oauth2Client, err := hydraAdmin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
		Body: &models.OAuth2Client{
			ClientID:     "integration-testing",
			ClientName:   "integration-testing",
			ClientSecret: "integration-testing",
			GrantTypes:   []string{"client_credentials"},
		},
		Context: ctx,
	})
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
		return
	}

	oidcProvider, err := oidc.NewProvider(ctx, "http://localhost:4444/")
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
		return
	}

	s.oauth2Cfg = clientcredentials.Config{
		ClientID:     oauth2Client.Payload.ClientID,
		ClientSecret: oauth2Client.Payload.ClientSecret,
		TokenURL:     oidcProvider.Endpoint().TokenURL,
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
