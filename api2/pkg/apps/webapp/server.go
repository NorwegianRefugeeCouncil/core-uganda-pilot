package webapp

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/nrc-no/core-kafka/pkg/auth"
	"github.com/nrc-no/core-kafka/pkg/sessionmanager"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	hydraAdminClient  *client.OryHydra
	hydraPublicClient *client.OryHydra
	renderFactory     *RendererFactory
	sessionManager    sessionmanager.Store
	credentialsClient *auth.CredentialsClient
	iam               iam.Interface
	cms               cms.Interface
	OauthClientID     string
	OauthClientSecret string
	router            *mux.Router
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
	credentialsClient *auth.CredentialsClient,
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
		Scope: "openid profile",
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

	h := &Server{
		hydraAdminClient:  hydraAdminClient,
		hydraPublicClient: hydraPublicClient,
		renderFactory:     renderFactory,
		sessionManager:    sm,
		credentialsClient: credentialsClient,
		iam:               iamClient,
		cms:               cmsClient,
		OauthClientID:     clientId,
		OauthClientSecret: clientSecret,
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
