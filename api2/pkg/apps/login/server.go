package login

import (
	"github.com/gorilla/mux"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
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
		ListenAddress: ":9001",
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
}

func NewServer(o *ServerOptions) (*Server, error) {

	hydraAdminURL, err := url.Parse(o.HydraAdminURL)
	if err != nil {
		return nil, err
	}

	cli := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     hydraAdminURL.Host,
		BasePath: hydraAdminURL.Path,
		Schemes:  []string{hydraAdminURL.Scheme},
	})

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

	collection := mongoClient.Database(o.MongoDatabase).Collection("credentials")

	srv := &Server{
		HydraAdmin: cli.Admin,
		Collection: collection,
		BCryptCost: o.BCryptCost,
	}

	router := mux.NewRouter()
	router.Path("/login").Methods("GET").HandlerFunc(srv.GetLogin)
	router.Path("/login").Methods("POST").HandlerFunc(srv.PostLogin)

	return srv, nil

}
