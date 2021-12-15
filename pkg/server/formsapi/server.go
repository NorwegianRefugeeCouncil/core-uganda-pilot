package formsapi

import (
	"context"
	"github.com/nrc-no/core/pkg/server/forms/handlers/database"
	"github.com/nrc-no/core/pkg/server/forms/handlers/folder"
	"github.com/nrc-no/core/pkg/server/forms/handlers/form"
	"github.com/nrc-no/core/pkg/server/forms/handlers/record"
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

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "public")
	if err != nil {
		return nil, err
	}

	container := genericServer.GoRestfulContainer

	databaseStore := store.NewDatabaseStore(options.StoreFactory)
	databaseHandler := database.NewHandler(databaseStore)
	container.Add(databaseHandler.WebService())

	folderStore := store.NewFolderStore(options.StoreFactory)
	folderHandler := folder.NewHandler(folderStore)
	container.Add(folderHandler.WebService())

	formStore := store.NewFormStore(options.StoreFactory)
	formHandler := form.NewHandler(formStore)
	container.Add(formHandler.WebService())

	recordStore := store.NewRecordStore(options.StoreFactory, formStore)
	recordHandler := record.NewHandler(recordStore, formStore)
	container.Add(recordHandler.WebService())

	s := &Server{
		options: options,
		Server:  genericServer,
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) {
	s.Server.Start(ctx)
}
