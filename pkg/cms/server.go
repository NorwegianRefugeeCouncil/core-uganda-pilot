package cms

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"path"
)

type Server struct {
	environment     string
	router          *mux.Router
	mongoClientFn   utils.MongoClientFn
	caseStore       *CaseStore
	caseTypeStore   *CaseTypeStore
	commentStore    *CommentStore
	HydraAdmin      admin.ClientService
	HydraHTTPClient *http.Client
}

func NewServerOrDie(ctx context.Context, o *server.GenericServerOptions) *Server {
	srv, err := NewServer(ctx, o)
	if err != nil {
		panic(err)
	}
	return srv
}

func NewServer(ctx context.Context, o *server.GenericServerOptions) (*Server, error) {
	caseStore, err := NewCaseStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	caseTypeStore, err := NewCaseTypeStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	commentStore, err := NewCommentStore(ctx, o.MongoClientFn, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		mongoClientFn:   o.MongoClientFn,
		environment:     o.Environment,
		caseStore:       caseStore,
		caseTypeStore:   caseTypeStore,
		commentStore:    commentStore,
		HydraAdmin:      o.HydraAdminClient.Admin,
		HydraHTTPClient: o.HydraHTTPClient,
	}

	router := mux.NewRouter()
	router.Use(srv.WithAuth())

	router.Path(server.CasesEndpoint).Methods("GET").HandlerFunc(srv.ListCases)
	router.Path(server.CasesEndpoint).Methods("POST").HandlerFunc(srv.PostCase)
	router.Path(path.Join(server.CasesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetCase)
	router.Path(path.Join(server.CasesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutCase)

	router.Path(server.CaseTypesEndpoint).Methods("GET").HandlerFunc(srv.ListCaseTypes)
	router.Path(server.CaseTypesEndpoint).Methods("POST").HandlerFunc(srv.postCaseType)
	router.Path(path.Join(server.CaseTypesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetCaseType)
	router.Path(path.Join(server.CaseTypesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.putCaseType)

	router.Path(server.CommentsEndpoint).Methods("GET").HandlerFunc(srv.ListComments)
	router.Path(server.CommentsEndpoint).Methods("POST").HandlerFunc(srv.PostComment)
	router.Path(path.Join(server.CommentsEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetComment)
	router.Path(path.Join(server.CommentsEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutComment)
	router.Path(path.Join(server.CommentsEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.DeleteComment)

	srv.router = router

	return srv, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
	utils.JSONResponse(w, status, data)
}

func (s *Server) getPathParam(param string, w http.ResponseWriter, req *http.Request, into *string) bool {
	return utils.GetPathParam(param, w, req, into)
}

func (s *Server) error(w http.ResponseWriter, err error) {
	utils.ErrorResponse(w, err)
}

func (s *Server) bind(req *http.Request, into interface{}) error {
	return utils.BindJSON(req, into)
}

func (s *Server) ResetDB(ctx context.Context, databaseName string) error {
	mongoClient, err := s.mongoClientFn(ctx)
	if err != nil {
		return err
	}

	if err := ClearCollections(ctx, mongoClient, databaseName); err != nil {
		return err
	}

	return nil
}

func ClearCollections(ctx context.Context, mongoCli *mongo.Client, databaseName string) error {
	return storage.ClearCollections(ctx, mongoCli, databaseName, AllCollections...)
}
