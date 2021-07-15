package webapp

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/login"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/ory/hydra-client-go/client/admin"
	"golang.org/x/oauth2"
	"net/http"
)

type Server struct {
	renderFactory   *RendererFactory
	sessionManager  sessionmanager.Store
	router          *mux.Router
	login           login.Interface
	HydraAdmin      admin.ClientService
	oidcVerifier    *oidc.IDTokenVerifier
	oauth2Config    *oauth2.Config
	environment     string
	iamAdminClient  iam.Interface
	baseURL         string
	iamScheme       string
	iamHost         string
	cmsScheme       string
	cmsHost         string
	HydraHTTPClient *http.Client
	IAMHTTPClient   *http.Client
	CMSHTTPClient   *http.Client
	Constants Constants
}

type ServerOptions struct {
	*server.GenericServerOptions
	TemplateDirectory string
	BaseURL           string
	IAMHost           string
	IAMScheme         string
	CMSHost           string
	CMSScheme         string
	AdminHTTPClient   *http.Client
	IDTokenVerifier   *oidc.IDTokenVerifier
	OAuth2Config      *oauth2.Config
	HydraHTTPClient   *http.Client
	IAMHTTPClient     *http.Client
	CMSHTTPClient     *http.Client
	LoginHTTPClient   *http.Client
}

func NewServer(options *ServerOptions) (*Server, error) {

	sm := sessionmanager.New(options.RedisPool, sessionmanager.Options{})

	renderFactory, err := NewRendererFactory(options.TemplateDirectory, sm)
	if err != nil {
		return nil, err
	}

	h := &Server{
		renderFactory:  renderFactory,
		sessionManager: sm,
		login: login.NewClientSet(&rest.RESTConfig{
			Scheme:     options.CMSScheme,
			Host:       options.CMSHost,
			HTTPClient: options.AdminHTTPClient,
		}),
		HydraAdmin:   options.HydraAdminClient.Admin,
		oidcVerifier: options.IDTokenVerifier,
		oauth2Config: options.OAuth2Config,
		environment:  options.Environment,
		iamAdminClient: iam.NewClientSet(&rest.RESTConfig{
			Scheme:     options.IAMScheme,
			Host:       options.IAMHost,
			HTTPClient: options.AdminHTTPClient,
		}),
		baseURL:         options.BaseURL,
		iamScheme:       options.IAMScheme,
		iamHost:         options.IAMHost,
		cmsScheme:       options.CMSScheme,
		cmsHost:         options.CMSHost,
		HydraHTTPClient: options.HydraHTTPClient,
		IAMHTTPClient:   options.IAMHTTPClient,
		CMSHTTPClient:   options.CMSHTTPClient,
		Constants: Constants{
			PartyDropdownLimit: 5,
		},
	}

	router := mux.NewRouter()
	router.Use(sm.LoadAndSave)
	router.Use(h.WithAuth())
	router.Path("/callback").HandlerFunc(h.Callback)
	router.Path("/login").HandlerFunc(h.Login)
	router.Path("/logout").HandlerFunc(h.Logout)
	router.Path("/").HandlerFunc(h.Individuals)
	router.Path("/individuals").HandlerFunc(h.Individuals)
	router.Path("/individuals/{id}").HandlerFunc(h.Individual)
	router.Path("/individuals/{id}/credentials").HandlerFunc(h.IndividualCredentials)
	router.Path("/teams").HandlerFunc(h.Teams)
	router.Path("/teams/pickparty").HandlerFunc(h.PickTeamParty)
	router.Path("/teams/{id}").HandlerFunc(h.Team)
	router.Path("/teams/{id}/invitemember").HandlerFunc(h.AddIndividualToTeam)
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
	router.Path("/comments").Methods("POST").HandlerFunc(h.PostComment)
	router.Path("/relationships/pickparty").HandlerFunc(h.PickRelationshipParty)

	h.router = router

	return h, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) IAMClient(ctx context.Context) iam.Interface {
	cfg := s.oauth2Config
	accessToken := s.sessionManager.GetString(ctx, "access-token")
	refreshToken := s.sessionManager.GetString(ctx, "refresh-token")
	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	httpClient := s.IAMHTTPClient
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	cli := cfg.Client(ctx, token)
	if len(accessToken) > 0 || len(refreshToken) > 0 {
		httpClient = cli
	}

	return iam.NewClientSet(&rest.RESTConfig{
		Scheme:     s.iamScheme,
		Host:       s.iamHost,
		HTTPClient: httpClient,
	})
}

func (s *Server) CMSClient(ctx context.Context) cms.Interface {
	cfg := s.oauth2Config
	accessToken := s.sessionManager.GetString(ctx, "access-token")
	refreshToken := s.sessionManager.GetString(ctx, "refresh-token")
	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	httpClient := s.CMSHTTPClient
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	cli := cfg.Client(ctx, token)
	if len(accessToken) > 0 || len(refreshToken) > 0 {
		httpClient = cli
	}
	return cms.NewClientSet(&rest.RESTConfig{
		Scheme:     s.cmsScheme,
		Host:       s.cmsHost,
		HTTPClient: httpClient,
	})
}
