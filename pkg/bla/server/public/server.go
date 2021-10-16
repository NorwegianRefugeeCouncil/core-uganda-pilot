package public

import (
	"github.com/nrc-no/core/pkg/bla/options"
	"github.com/nrc-no/core/pkg/bla/server/generic"
	"github.com/nrc-no/core/pkg/bla/server/public/handlers/database"
	"github.com/nrc-no/core/pkg/bla/server/public/handlers/folder"
	"github.com/nrc-no/core/pkg/bla/server/public/handlers/form"
	"github.com/nrc-no/core/pkg/bla/server/public/handlers/record"
	"github.com/nrc-no/core/pkg/bla/store"
)

type Server struct {
	*generic.Server
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

	container := genericServer.Container

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
	recordHandler := record.NewHandler(recordStore)
	container.Add(recordHandler.WebService())

	s := &Server{
		Server: genericServer,
	}

	return s, nil
}
