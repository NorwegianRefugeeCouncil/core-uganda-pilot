package iam

import (
	"github.com/gorilla/mux"
	"net/http"
)

type ServerOptions struct {
	Address string
}

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	router := mux.NewRouter()
	srv := &Server{
		router: router,
	}
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
