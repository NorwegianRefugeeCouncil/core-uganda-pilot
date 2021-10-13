package seeder

import (
	"context"
	"github.com/nrc-no/core/pkg/generic/server"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

type Server struct {
	environment     string
	context 		context.Context
	mongoClientFn   utils.MongoClientFn
	mongoDatabase 	string
}

func NewServer(ctx context.Context, o *server.GenericServerOptions) (*Server, error) {

	srv := &Server{
		environment:     o.Environment,
		context: 		 ctx,
		mongoClientFn:   o.MongoClientFn,
		mongoDatabase:   o.MongoDatabase,
	}

	return srv, nil

}

func (s *Server) ClearDB() {
	if err := Clear(s.context, s.mongoClientFn, s.mongoDatabase); err != nil {
		panic(err)
	}
}

func (s *Server) SeedDB() {
	if err := Seed(s.context, s.mongoClientFn, s.mongoDatabase); err != nil {
		panic(err)
	}
}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
	utils.JSONResponse(w, status, data)
}

func (s *Server) error(w http.ResponseWriter, err error) {
	utils.ErrorResponse(w, err)
}

func (s *Server) bind(req *http.Request, into interface{}) error {
	return utils.BindJSON(req, into)
}


