package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
)

type ServerOptions struct {
	ListenAddress string
	MongoHosts    []string
	MongoDatabase string
	MongoUsername string
	MongoPassword string
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		ListenAddress: ":9001",
		MongoHosts:    []string{"mongo://localhost:27017"},
	}
}

func (o *ServerOptions) WithMongoHosts(hosts []string) *ServerOptions {
	o.MongoHosts = hosts
	return o
}
func (o *ServerOptions) WithMongoDatabase(mongoDatabase string) *ServerOptions {
	o.MongoDatabase = mongoDatabase
	return o
}
func (o *ServerOptions) WithMongoUsername(mongoUsername string) *ServerOptions {
	o.MongoUsername = mongoUsername
	return o
}
func (o *ServerOptions) WithMongoPassword(mongoPassword string) *ServerOptions {
	o.MongoPassword = mongoPassword
	return o
}
func (o *ServerOptions) WithListenAddress(address string) *ServerOptions {
	o.ListenAddress = address
	return o
}

func (o *ServerOptions) Flags(fs pflag.FlagSet) {
	fs.StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "Server listen address")
	fs.StringSliceVar(&o.MongoHosts, "mongo-url", o.MongoHosts, "Mongo url")
	fs.StringVar(&o.MongoDatabase, "mongo-database", o.MongoDatabase, "Mongo database")
	fs.StringVar(&o.MongoUsername, "mongo-username", o.MongoUsername, "Mongo username")
	fs.StringVar(&o.MongoPassword, "mongo-password", o.MongoPassword, "Mongo password")
}

type Server struct {
	router        *mux.Router
	mongoClient   *mongo.Client
	caseStore     *CaseStore
	caseTypeStore *CaseTypeStore
}

func NewServer(ctx context.Context, o *ServerOptions) (*Server, error) {
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

	router := mux.NewRouter()

	caseStore, err := NewCaseStore(ctx, mongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	caseTypeStore, err := NewCaseTypeStore(ctx, mongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		router:        router,
		mongoClient:   mongoClient,
		caseStore:     caseStore,
		caseTypeStore: caseTypeStore,
	}

	router.Use(srv.WithAuth(ctx))

	router.Path("/apis/cms/v1/cases").Methods("GET").HandlerFunc(srv.ListCases)
	router.Path("/apis/cms/v1/cases").Methods("POST").HandlerFunc(srv.PostCase)
	router.Path("/apis/cms/v1/cases/{id}").Methods("GET").HandlerFunc(srv.GetCase)
	router.Path("/apis/cms/v1/cases/{id}").Methods("PUT").HandlerFunc(srv.PutCase)

	router.Path("/apis/cms/v1/casetypes").Methods("GET").HandlerFunc(srv.ListCaseTypes)
	router.Path("/apis/cms/v1/casetypes").Methods("POST").HandlerFunc(srv.PostCaseType)
	router.Path("/apis/cms/v1/casetypes/{id}").Methods("GET").HandlerFunc(srv.GetCaseType)
	router.Path("/apis/cms/v1/casetypes/{id}").Methods("PUT").HandlerFunc(srv.PutCaseType)

	return srv, nil

}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) JSON(w http.ResponseWriter, status int, data interface{}) {
	responseBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

func (s *Server) GetPathParam(param string, w http.ResponseWriter, req *http.Request, into *string) bool {
	id, ok := mux.Vars(req)[param]
	if !ok || len(id) == 0 {
		err := fmt.Errorf("path parameter '%s' not found in path", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	*into = id
	return true
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
