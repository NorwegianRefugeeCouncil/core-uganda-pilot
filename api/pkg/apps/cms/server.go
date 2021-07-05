package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/ory/hydra-client-go/client/admin"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
)

type Server struct {
	environment     string
	router          *mux.Router
	mongoClient     *mongo.Client
	caseStore       *CaseStore
	caseTypeStore   *CaseTypeStore
	commentStore    *CommentStore
	HydraAdmin      admin.ClientService
	HydraHttpClient *http.Client
}

func NewServer(ctx context.Context, o *server.GenericServerOptions) (*Server, error) {

	caseStore, err := NewCaseStore(ctx, o.MongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	caseTypeStore, err := NewCaseTypeStore(ctx, o.MongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	commentStore, err := NewCommentStore(ctx, o.MongoClient, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		mongoClient:     o.MongoClient,
		environment:     o.Environment,
		caseStore:       caseStore,
		caseTypeStore:   caseTypeStore,
		commentStore:    commentStore,
		HydraAdmin:      o.HydraAdminClient.Admin,
		HydraHttpClient: o.HydraHTTPClient,
	}

	router := mux.NewRouter()
	router.Use(srv.WithAuth())

	casesEP := server.Endpoints["cases"]
	caseTypesEP := server.Endpoints["casetypes"]
	commentsEP := server.Endpoints["comments"]

	router.Path(casesEP).Methods("GET").HandlerFunc(srv.ListCases)
	router.Path(casesEP).Methods("POST").HandlerFunc(srv.PostCase)
	router.Path(casesEP + "/{id}").Methods("GET").HandlerFunc(srv.GetCase)
	router.Path(casesEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutCase)

	router.Path(caseTypesEP).Methods("GET").HandlerFunc(srv.ListCaseTypes)
	router.Path(caseTypesEP).Methods("POST").HandlerFunc(srv.PostCaseType)
	router.Path(caseTypesEP + "/{id}").Methods("GET").HandlerFunc(srv.GetCaseType)
	router.Path(caseTypesEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutCaseType)

	router.Path(commentsEP).Methods("GET").HandlerFunc(srv.ListComments)
	router.Path(commentsEP).Methods("POST").HandlerFunc(srv.PostComment)
	router.Path(commentsEP + "/{id}").Methods("GET").HandlerFunc(srv.GetComment)
	router.Path(commentsEP + "/{id}").Methods("PUT").HandlerFunc(srv.PutComment)
	router.Path(commentsEP + "/{id}").Methods("PUT").HandlerFunc(srv.DeleteComment)

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
