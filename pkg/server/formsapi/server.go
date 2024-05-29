package formsapi

import (
	"context"
	"github.com/nrc-no/core/pkg/server/formsapi/handlers/database"
	"github.com/nrc-no/core/pkg/server/formsapi/handlers/folder"
	"github.com/nrc-no/core/pkg/server/formsapi/handlers/form"
	"github.com/nrc-no/core/pkg/server/formsapi/handlers/record"
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	client2 "github.com/nrc-no/core/pkg/zanzibar"
)

type Server struct {
	*generic.Server
	options Options
}

type Options struct {
	options.ServerOptions
	StoreFactory store.Factory
	ZanzibarClient client2.ZanzibarClient
}

func NewServer(options Options) (*Server, error) {

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "forms-api")
	if err != nil {
		return nil, err
	}

	container := genericServer.GoRestfulContainer

	databaseStore := store.NewDatabaseStore(options.StoreFactory)
	databaseHandler := database.NewHandler(databaseStore, options.ZanzibarClient)
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
