package login

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/ory/hydra-client-go/client/admin"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"io/ioutil"
	"net/http"
)

type ServerOptions struct {
	*server.GenericServerOptions
	BCryptCost      int
	AdminHTTPClient *http.Client
	IAMHost         string
	IAMScheme       string
}

type Server struct {
	HydraAdmin admin.ClientService
	Collection *mongo.Collection
	BCryptCost int
	router     *mux.Router
	template   *template.Template
	iam        iam.Interface
}

func NewServer(ctx context.Context, o *ServerOptions) (*Server, error) {

	iamCli := iam.NewClientSet(&rest.RESTConfig{
		Scheme:     o.IAMScheme,
		Host:       o.IAMHost,
		HTTPClient: o.AdminHTTPClient,
	})

	collection := o.MongoClient.Database(o.MongoDatabase).Collection("credentials")

	srv := &Server{
		HydraAdmin: o.HydraAdminClient.Admin,
		Collection: collection,
		BCryptCost: o.BCryptCost,
		iam:        iamCli,
	}

	router := mux.NewRouter()
	router.Path("/auth/logout").Methods("GET").HandlerFunc(srv.GetLogoutForm)
	router.Path("/auth/login").Methods("GET").HandlerFunc(srv.GetLoginForm)
	router.Path("/auth/login").Methods("POST").HandlerFunc(srv.PostLoginForm)
	router.Path("/auth/consent").Methods("GET").HandlerFunc(srv.GetConsent)
	router.Path("/auth/consent").Methods("POST").HandlerFunc(srv.PostConsent)
	router.Path("/apis/login/v1/credentials").
		Methods("POST").
		Handler(srv.WithAuth()(http.HandlerFunc(srv.PostLoginForm)))

	srv.router = router

	tpl, err := template.ParseGlob("pkg/apps/login/templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	srv.template = tpl

	return srv, nil

}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) Error(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) Bind(req *http.Request, into interface{}) error {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bodyBytes, &into); err != nil {
		return err
	}
	return nil
}
