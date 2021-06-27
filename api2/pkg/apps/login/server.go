package login

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/nrc-no/core-kafka/pkg/rest"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2/clientcredentials"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ServerOptions struct {
	ListenAddress string
	HydraAdminURL string
	MongoHosts    []string
	MongoDatabase string
	MongoUsername string
	MongoPassword string
	BCryptCost    int
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		HydraAdminURL: "http://localhost:4445",
		ListenAddress: ":9000",
		MongoHosts:    []string{"mongo://localhost:27017"},
		BCryptCost:    15,
	}
}

func (o *ServerOptions) Flags(fs pflag.FlagSet) {
	fs.StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "Server listen address")
	fs.StringVar(&o.HydraAdminURL, "hydra-admin-url", o.HydraAdminURL, "Ory Hydra admin URL")
	fs.StringSliceVar(&o.MongoHosts, "mongo-url", o.MongoHosts, "Mongo url")
	fs.StringVar(&o.MongoDatabase, "mongo-database", o.MongoDatabase, "Mongo database")
	fs.StringVar(&o.MongoUsername, "mongo-username", o.MongoUsername, "Mongo username")
	fs.StringVar(&o.MongoPassword, "mongo-password", o.MongoPassword, "Mongo password")
	fs.IntVar(&o.BCryptCost, "bcrypt-cost", o.BCryptCost, "BCrypt cost parameter")
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

	hydraAdminURL, err := url.Parse(o.HydraAdminURL)
	if err != nil {
		return nil, err
	}

	hydraClient := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     hydraAdminURL.Host,
		BasePath: hydraAdminURL.Path,
		Schemes:  []string{hydraAdminURL.Scheme},
	})
	hydraAdmin := hydraClient.Admin

	cli := &models.OAuth2Client{
		ClientID:     "login",
		ClientName:   "Login",
		ClientSecret: "somesupersecret",
		GrantTypes: []string{
			"client_credentials",
		},
		ResponseTypes:           []string{"token", "refresh_token"},
		TokenEndpointAuthMethod: "client_secret_post",
	}

	if err := createOauthClient(ctx, hydraAdmin, cli); err != nil {
		return nil, err
	}

	clientCredsCfg := clientcredentials.Config{
		ClientID:     cli.ClientID,
		ClientSecret: cli.ClientSecret,
		TokenURL:     "http://localhost:4444/oauth2/token",
	}
	adminCli := clientCredsCfg.Client(ctx)

	iamCli := iam.NewClientSet(&rest.RESTConfig{
		Scheme:     "http",
		Host:       "localhost:9000",
		HTTPClient: adminCli,
	})

	go func() {
		time.Sleep(2 * time.Second)
		_, err := iamCli.Staff().List(ctx, iam.StaffListOptions{})
		if err != nil {
			logrus.WithError(err).Errorf("")
		}
	}()

	mongoClient, err := mongo.NewClient(
		options.Client().
			SetHosts(o.MongoHosts).
			SetAuth(options.Credential{
				Username: o.MongoUsername,
				Password: o.MongoPassword,
			}))

	if err != nil {
		return nil, err
	}

	if err := mongoClient.Connect(ctx); err != nil {
		return nil, err
	}

	collection := mongoClient.Database(o.MongoDatabase).Collection("credentials")

	srv := &Server{
		HydraAdmin: hydraAdmin,
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

func createOauthClient(
	ctx context.Context,
	hydraAdmin admin.ClientService,
	cli *models.OAuth2Client) error {
	_, err := hydraAdmin.CreateOAuth2Client(&admin.CreateOAuth2ClientParams{
		Body:    cli,
		Context: ctx,
	})
	if err != nil {
		if strings.Contains(err.Error(), "createOAuth2ClientConflict") {
			_, err = hydraAdmin.UpdateOAuth2Client(&admin.UpdateOAuth2ClientParams{
				ID:      cli.ClientID,
				Body:    cli,
				Context: ctx,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
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
