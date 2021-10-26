package attachments

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/storage"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/ory/hydra-client-go/client/admin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"path"
)

type Server struct {
	environment     string
	router         *mux.Router
	mongoClientSrc storage.MongoClientSrc
	store          *AttachmentStore
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
	attachmentStore, err := NewAttachmentStore(ctx, o.MongoClientSrc, o.MongoDatabase)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		environment:     o.Environment,
		mongoClientSrc:  o.MongoClientSrc,
		store:           attachmentStore,
		HydraAdmin:      o.HydraAdminClient.Admin,
		HydraHTTPClient: o.HydraHTTPClient,
	}

	router := mux.NewRouter()
	router.Use(srv.WithAuth())

	router.Path(server.AttachmentsEndpoint).Methods(http.MethodGet).HandlerFunc(srv.ListAttachments)
	router.Path(server.AttachmentsEndpoint).Methods(http.MethodPost).HandlerFunc(srv.PostAttachment)
	router.Path(path.Join(server.AttachmentsEndpoint, "{id}")).Methods(http.MethodGet).HandlerFunc(srv.GetAttachment)
	router.Path(path.Join(server.AttachmentsEndpoint, "{id}")).Methods(http.MethodPut).HandlerFunc(srv.PutAttachment)

	srv.router = router

	return srv, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) Error(w http.ResponseWriter, err error) {
	utils.ErrorResponse(w, err)
}

func (s *Server) Bind(req *http.Request, into interface{}) error {
	return utils.BindJSON(req, into)
}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
	utils.JSONResponse(w, status, data)
}

func (s *Server) GetPathParam(param string, w http.ResponseWriter, req *http.Request, into *string) bool {
	return utils.GetPathParam(param, w, req, into)
}

func (s *Server) ResetDB(ctx context.Context, databaseName string) error {
	mongoClient, err := s.mongoClientSrc.GetMongoClient()
	if err != nil {
		return err
	}
	// Delete attachments
	_, err = mongoClient.Database(databaseName).Collection("attachments").DeleteMany(ctx, bson.M{})
	return err
}
