package webapp

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/nrc-no/core-kafka/pkg/apps/login"
	"github.com/nrc-no/core-kafka/pkg/rest"
	"github.com/nrc-no/core-kafka/pkg/sessionmanager"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	hydraAdminClient  *client.OryHydra
	hydraPublicClient *client.OryHydra
	renderFactory     *RendererFactory
	sessionManager    sessionmanager.Store
	iam               iam.Interface
	cms               cms.Interface
	OauthClientID     string
	OauthClientSecret string
	router            *mux.Router
	login             login.Interface
	HydraAdmin        admin.ClientService
}

type ServerOptions struct {
	RedisNetwork            string
	RedisAddress            string
	RedisMaxIdleConnections int
	TemplateDirectory       string
}

func NewServer(
	options ServerOptions,
	hydraAdminClient *client.OryHydra,
	hydraPublicClient *client.OryHydra,
	iamClient *iam.ClientSet,
	cmsClient *cms.ClientSet,
) (*Server, error) {

	ctx := context.Background()

	pool := &redis.Pool{
		MaxIdle: options.RedisMaxIdleConnections,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(options.RedisNetwork, options.RedisAddress)
		},
	}

	sm := sessionmanager.New(pool, sessionmanager.Options{})

	renderFactory, err := NewRendererFactory(options.TemplateDirectory)
	if err != nil {
		return nil, err
	}

	e, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fmt.Println(e)

	clientId := "webapp"
	clientSecret := "somesupersecret"

	cli := &models.OAuth2Client{
		ClientID:     clientId,
		ClientName:   "Web App",
		ClientSecret: clientSecret,
		GrantTypes: []string{
			"authorization_code",
			"id_token",
			"access_token",
			"refresh_token",
		},
		RedirectUris: []string{
			"http://localhost:9000/callback",
		},
		ResponseTypes: []string{
			"token",
			"code",
		},
		Scope:                   "openid profile",
		TokenEndpointAuthMethod: "client_secret_post",
	}

	_, err = hydraAdminClient.Admin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
		Body:    cli,
		Context: ctx,
	})
	if err != nil {
		if strings.Contains(err.Error(), "createOAuth2ClientConflict") {
			_, err = hydraAdminClient.Admin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
				Body:    cli,
				ID:      cli.ClientID,
				Context: ctx,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	oauth2Cfg := clientcredentials.Config{
		ClientID:     cli.ClientID,
		ClientSecret: cli.ClientSecret,
		TokenURL:     "http://localhost:4444/oauth2/token",
	}
	loginClient := login.NewClientSet(&rest.RESTConfig{
		Scheme:     "http",
		Host:       "localhost:9000",
		HTTPClient: oauth2Cfg.Client(ctx),
	})

	h := &Server{
		hydraAdminClient:  hydraAdminClient,
		hydraPublicClient: hydraPublicClient,
		renderFactory:     renderFactory,
		sessionManager:    sm,
		iam:               iamClient,
		cms:               cmsClient,
		login:             loginClient,
		OauthClientID:     clientId,
		OauthClientSecret: clientSecret,
		HydraAdmin:        hydraAdminClient.Admin,
	}

	router := mux.NewRouter()
	router.Use(sm.LoadAndSave)
	router.Use(h.WithAuth(ctx))
	router.Path("/callback").HandlerFunc(h.Callback)
	router.Path("/login").HandlerFunc(h.Login)
	router.Path("/").HandlerFunc(h.Individuals)
	router.Path("/individuals").HandlerFunc(h.Individuals)
	router.Path("/individuals/{id}").HandlerFunc(h.Individual)
	router.Path("/individuals/{id}/credentials").HandlerFunc(h.IndividualCredentials)
	router.Path("/teams").HandlerFunc(h.Teams)
	router.Path("/teams/{id}").HandlerFunc(h.Team)
	router.Path("/cases").HandlerFunc(h.Cases)
	router.Path("/cases/new").HandlerFunc(h.NewCase)
	router.Path("/cases/{id}").HandlerFunc(h.Case)
	router.Path("/settings").HandlerFunc(h.Settings)
	router.Path("/settings/attributes").HandlerFunc(h.Attributes)
	router.Path("/settings/attributes/new").HandlerFunc(h.NewAttribute)
	router.Path("/settings/attributes/{id}").HandlerFunc(h.Attribute)
	router.Path("/settings/relationshiptypes").HandlerFunc(h.RelationshipTypes)
	router.Path("/settings/relationshiptypes/new").HandlerFunc(h.NewRelationshipType)
	router.Path("/settings/relationshiptypes/{id}").HandlerFunc(h.RelationshipType)
	router.Path("/settings/partytypes").HandlerFunc(h.PartyTypes)
	router.Path("/settings/partytypes/{id}").HandlerFunc(h.PartyType)
	router.Path("/settings/casetypes").HandlerFunc(h.CaseTypes)
	router.Path("/settings/casetypes/new").HandlerFunc(h.NewCaseType)
	router.Path("/settings/casetypes/{id}").HandlerFunc(h.CaseType)
	router.Path("/settings/authclients").HandlerFunc(h.AuthClients)
	router.Path("/settings/authclients/{id}").HandlerFunc(h.AuthClient)
	router.Path("/settings/authclients/{id}/newsecret").HandlerFunc(h.AuthClientNewSecret)
	router.Path("/settings/authclients/{id}/delete").HandlerFunc(h.DeleteAuthClient)

	h.router = router

	return h, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) IAMClient(ctx context.Context) iam.Interface {
	cfg := oauth2.Config{
		ClientID:     s.OauthClientID,
		ClientSecret: s.OauthClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:4444/oauth2/auth",
			TokenURL: "http://localhost:4444/oauth2/token",
		},
	}
	token := &oauth2.Token{
		AccessToken:  s.sessionManager.GetString(ctx, "access-token"),
		RefreshToken: s.sessionManager.GetString(ctx, "refresh-token"),
	}
	cli := cfg.Client(ctx, token)
	return iam.NewClientSet(&rest.RESTConfig{
		Scheme:     "http",
		Host:       "localhost:9000",
		HTTPClient: cli,
	})
}

func (s *Server) CMSClient(ctx context.Context) cms.Interface {
	cfg := oauth2.Config{
		ClientID:     s.OauthClientID,
		ClientSecret: s.OauthClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:4444/oauth2/auth",
			TokenURL: "http://localhost:4444/oauth2/token",
		},
	}
	token := &oauth2.Token{
		AccessToken:  s.sessionManager.GetString(ctx, "access-token"),
		RefreshToken: s.sessionManager.GetString(ctx, "refresh-token"),
	}
	cli := cfg.Client(ctx, token)
	return cms.NewClientSet(&rest.RESTConfig{
		Scheme:     "http",
		Host:       "localhost:9000",
		HTTPClient: cli,
	})
}
