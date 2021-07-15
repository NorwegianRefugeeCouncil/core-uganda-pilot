package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"io/ioutil"
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
	HydraHttpClient *http.Client
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
		HydraHttpClient: o.HydraHTTPClient,
	}

	router := mux.NewRouter()
	router.Use(srv.WithAuth())

	router.Path(server.CasesEndpoint).Methods("GET").HandlerFunc(srv.ListCases)
	router.Path(server.CasesEndpoint).Methods("POST").HandlerFunc(srv.PostCase)
	router.Path(path.Join(server.CasesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetCase)
	router.Path(path.Join(server.CasesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutCase)

	router.Path(server.CaseTypesEndpoint).Methods("GET").HandlerFunc(srv.ListCaseTypes)
	router.Path(server.CaseTypesEndpoint).Methods("POST").HandlerFunc(srv.PostCaseType)
	router.Path(path.Join(server.CaseTypesEndpoint, "{id}")).Methods("GET").HandlerFunc(srv.GetCaseType)
	router.Path(path.Join(server.CaseTypesEndpoint, "{id}")).Methods("PUT").HandlerFunc(srv.PutCaseType)

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

func (s *Server) JSON(w http.ResponseWriter, status int, data interface{}) {
	responseBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseBytes)
	if err != nil {
		return
	}
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
