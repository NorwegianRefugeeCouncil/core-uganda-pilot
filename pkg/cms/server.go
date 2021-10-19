package cms

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"go.mongodb.org/mongo-driver/bson"
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
	router.Path(server.CasesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.ListCases)
	router.Path(server.CasesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.PostCase)
	router.Path(path.Join(server.CasesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.GetCase)
	router.Path(path.Join(server.CasesEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.PutCase)

	router.Path(server.CaseTypesEndpoint).Methods(http.MethodGet).HandlerFunc(srv.ListCaseTypes)
	router.Path(server.CaseTypesEndpoint).Methods(http.MethodPost).HandlerFunc(srv.postCaseType)
	router.Path(path.Join(server.CaseTypesEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.GetCaseType)
	router.Path(path.Join(server.CaseTypesEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.putCaseType)

	router.Path(server.CommentsEndpoint).Methods(http.MethodGet).HandlerFunc(srv.ListComments)
	router.Path(server.CommentsEndpoint).Methods(http.MethodPost).HandlerFunc(srv.PostComment)
	router.Path(path.Join(server.CommentsEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.GetComment)
	router.Path(path.Join(server.CommentsEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.PutComment)
	router.Path(path.Join(server.CommentsEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.DeleteComment)

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
	// Delete CaseTypes
	_, err = mongoClient.Database(databaseName).Collection("caseTypes").DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}
	// Delete Cases
	_, err = mongoClient.Database(databaseName).Collection("cases").DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}
	// Delete Comments
	_, err = mongoClient.Database(databaseName).Collection("comments").DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}
	return mongoClient.Database(databaseName).Drop(ctx)
}
