package formsapi

import (
	"context"

	"github.com/nrc-no/core/pkg/common"
	entityController "github.com/nrc-no/core/pkg/server/core-db/controllers/entity"
	entityModel "github.com/nrc-no/core/pkg/server/core-db/models/entity"
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
)

type Server struct {
	*generic.Server
	options Options
}

type Options struct {
	options.ServerOptions
	StoreFactory store.Factory
}

func NewServer(options Options) (*Server, error) {

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "core-db-api")
	if err != nil {
		return nil, err
	}

	container := genericServer.GoRestfulContainer

	transactionManager := common.NewTransactionManager(options.StoreFactory)

	entityModel := entityModel.NewEntityPostgresModel(options.StoreFactory)
	entityController := entityController.NewController(entityModel, transactionManager)
	container.Add(entityController.WebService())

	s := &Server{
		options: options,
		Server:  genericServer,
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) {
	s.Server.Start(ctx)
}
