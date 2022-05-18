package data

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nrc-no/core/pkg/server/data/api"
	"github.com/nrc-no/core/pkg/server/data/engine"
	"github.com/nrc-no/core/pkg/server/data/handler"
	"github.com/nrc-no/core/pkg/server/data/utils"
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/options"
)

type Server struct {
	*generic.Server
	options Options
}

type Options struct {
	options.ServerOptions
}

func NewServer(options Options) (*Server, error) {
	ctx := context.Background()

	// create the generic server
	genericServer, err := generic.NewGenericServer(options.ServerOptions, "data")
	if err != nil {
		return nil, err
	}

	// create the database connection
	db, err := sqlx.ConnectContext(ctx, "sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// create the engine
	e, err := engine.NewEngine(
		ctx,
		func(ctx context.Context) (api.Transaction, error) {
			return utils.NewTransaction(ctx, db, nil)
		},
		&utils.UUIDGenerator{},
		&utils.Md5RevGenerator{},
		&utils.Clock{},
		"sqlite",
	)
	if err != nil {
		return nil, err
	}

	// create the handlers
	h := handler.NewHandler(e)
	genericServer.GoRestfulContainer.Add(h.WebService())

	// create the server
	s := &Server{
		options: options,
		Server:  genericServer,
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) {
	s.Server.Start(ctx)
}
