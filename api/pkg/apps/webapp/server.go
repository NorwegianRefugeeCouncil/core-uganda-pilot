package webapp

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/cms"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/apps/login"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

type Server struct {
	renderFactory       *RendererFactory
	sessionManager      sessionmanager.Store
	router              *mux.Router
	login               login.Interface
	HydraAdmin          admin.ClientService
	oidcVerifier        *oidc.IDTokenVerifier
	privateOauth2Config *oauth2.Config
	environment         string
	iamAdminClient      iam.Interface
	baseURL             string
	iamScheme           string
	iamHost             string
	cmsScheme           string
	cmsHost             string
	HydraHTTPClient     *http.Client
	IAMHTTPClient       *http.Client
	CMSHTTPClient       *http.Client
	Constants           Constants
	publicOauth2Config  *oauth2.Config
}

type ServerOptions struct {
	*server.GenericServerOptions
	TemplateDirectory   string
	BaseURL             string
	IAMHost             string
	IAMScheme           string
	CMSHost             string
	CMSScheme           string
	AdminHTTPClient     *http.Client
	IDTokenVerifier     *oidc.IDTokenVerifier
	PrivateOAuth2Config *oauth2.Config
	HydraHTTPClient     *http.Client
	IAMHTTPClient       *http.Client
	CMSHTTPClient       *http.Client
	LoginHTTPClient     *http.Client
	PublicOauth2Config  *oauth2.Config
}

func NewServer(options *ServerOptions) (*Server, error) {

	h := &Server{
		login: login.NewClientSet(&rest.RESTConfig{
			Scheme:     options.CMSScheme,
			Host:       options.CMSHost,
			HTTPClient: options.AdminHTTPClient,
		}),
		HydraAdmin:          options.HydraAdminClient.Admin,
		oidcVerifier:        options.IDTokenVerifier,
		privateOauth2Config: options.PrivateOAuth2Config,
		publicOauth2Config:  options.PublicOauth2Config,
		environment:         options.Environment,
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

	sm, err := sessionmanager.New(options.RedisStore)
	if err != nil {
		logrus.WithError(err).Errorf("failed to create session manager")
		return nil, err
	}

	renderFactory, err := NewRendererFactory(options.TemplateDirectory, sm)
	if err != nil {
		logrus.WithError(err).Errorf("failed to create renderer")
		return nil, err
	}

	h.sessionManager = sm
	h.renderFactory = renderFactory

	router := mux.NewRouter()
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
	router.Path("/cases").Methods("GET").HandlerFunc(h.Cases)
	router.Path("/cases/new").Methods("GET").HandlerFunc(h.NewCase)
	router.Path("/cases/{id}").Methods("GET").HandlerFunc(h.Case)
	router.Path("/cases").Methods("POST").HandlerFunc(h.PostCase)
	router.Path("/cases/{id}").Methods("POST").HandlerFunc(h.PostCase)
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
	router.Path("/comments").Methods("POST").HandlerFunc(h.PostComment)
	router.Path("/relationships/pickparty").HandlerFunc(h.PickRelationshipParty)

	h.router = router

	return h, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) IAMClient(req *http.Request) (iam.Interface, error) {
	httpClient, err := utils.GetOauth2HttpClient(s.sessionManager, req, s.privateOauth2Config, s.IAMHTTPClient)
	if err != nil {
		return nil, err
	}
	return iam.NewClientSet(&rest.RESTConfig{
		Scheme:     s.iamScheme,
		Host:       s.iamHost,
		HTTPClient: httpClient,
	}), nil
}

func (s *Server) CMSClient(req *http.Request) (cms.Interface, error) {
	httpClient, err := utils.GetOauth2HttpClient(s.sessionManager, req, s.privateOauth2Config, s.CMSHTTPClient)
	if err != nil {
		return nil, err
	}
	return cms.NewClientSet(&rest.RESTConfig{
		Scheme:     s.cmsScheme,
		Host:       s.cmsHost,
		HTTPClient: httpClient,
	}), nil
}

func (s *Server) Error(w http.ResponseWriter, err error) {
	utils.ErrorResponse(w, err)
}
func (s *Server) GetPathParam(param string, w http.ResponseWriter, req *http.Request, into *string) bool {
	return utils.GetPathParam(param, w, req, into)
}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
	utils.JSONResponse(w, status, data)
}
