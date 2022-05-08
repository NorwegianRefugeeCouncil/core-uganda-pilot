package formsapi

import (
	"context"

	"github.com/nrc-no/core/pkg/common"
	entityController "github.com/nrc-no/core/pkg/server/core-db/controllers/entity"
	entityStore "github.com/nrc-no/core/pkg/server/core-db/stores/entity"
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

	transactionStore := common.NewTransactionStore(options.StoreFactory)

	entityStore := entityStore.NewEntityPostgresStore(options.StoreFactory)
	entityController := entityController.NewController(entityStore, transactionStore)
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
