package login

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/iam"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"path"
)

type ServerOptions struct {
	*server.GenericServerOptions
	BCryptCost        int
	AdminHTTPClient   *http.Client
	IAMHost           string
	IAMScheme         string
	TemplateDirectory string
}

type Server struct {
	HydraAdmin              admin.ClientService
	BCryptCost              int
	router                  *mux.Router
	template                *template.Template
	iam                     iam.Interface
	HydraHTTPClient         *http.Client
	credentialsCollectionFn func() (*mongo.Collection, error)
}

func NewServer(ctx context.Context, o *ServerOptions) (*Server, error) {
	iamCli := iam.NewClientSet(&rest.Config{
		Scheme:     o.IAMScheme,
		Host:       o.IAMHost,
		HTTPClient: o.AdminHTTPClient,
	})

	srv := &Server{
		HydraAdmin:      o.HydraAdminClient.Admin,
		BCryptCost:      o.BCryptCost,
		iam:             iamCli,
		HydraHTTPClient: o.HydraHTTPClient,
		credentialsCollectionFn: func() (*mongo.Collection, error) {
			mongoClient, err := o.MongoClientSrc.GetMongoClient()
			if err != nil {
				logrus.WithError(err).Errorf("unable to get mongo client")
				return nil, err
			}
			collection := mongoClient.Database(o.MongoDatabase).Collection(CredentialsCollection)
			return collection, nil
		},
	}

	router := mux.NewRouter()
	router.Path("/auth/logout").Methods(http.MethodGet).HandlerFunc(srv.GetLogoutForm)
	router.Path("/auth/login").Methods(http.MethodGet).HandlerFunc(srv.GetLoginForm)
	router.Path("/auth/login").Methods(http.MethodPost).HandlerFunc(srv.PostLoginForm)
	router.Path("/auth/consent").Methods(http.MethodGet).HandlerFunc(srv.GetConsent)
	router.Path("/auth/consent").Methods(http.MethodPost).HandlerFunc(srv.PostConsent)
	router.Path("/apis/login/v1/credentials").
		Methods(http.MethodPost).
		Handler(srv.WithAuth()(http.HandlerFunc(srv.PostCredentials)))

	srv.router = router

	tpl, err := template.ParseGlob(path.Join(o.TemplateDirectory, "*.gohtml"))
	if err != nil {
		logrus.WithError(err).Errorf("failed to parse templates")
		return nil, err
	}

	srv.template = tpl

	return srv, nil
}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
	utils.JSONResponse(w, status, data)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) Error(w http.ResponseWriter, err error) {
	logrus.WithError(err).Error()
	utils.ErrorResponse(w, err)
}

func (s *Server) Bind(req *http.Request, into interface{}) error {
	return utils.BindJSON(req, into)
}

func ClearCollections(ctx context.Context, mongoCli *mongo.Client, databaseName string) error {
	return storage.ClearCollections(ctx, mongoCli, databaseName, AllCollections...)
}
