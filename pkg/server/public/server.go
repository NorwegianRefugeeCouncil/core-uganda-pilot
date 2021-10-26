package public

import (
	"github.com/nrc-no/core/pkg/options"
	"github.com/nrc-no/core/pkg/server/generic"
	"github.com/nrc-no/core/pkg/server/public/handlers/database"
	"github.com/nrc-no/core/pkg/server/public/handlers/folder"
	"github.com/nrc-no/core/pkg/server/public/handlers/form"
	"github.com/nrc-no/core/pkg/server/public/handlers/record"
	store2 "github.com/nrc-no/core/pkg/store"
)

type Server struct {
	*generic.Server
}

type Options struct {
	options.ServerOptions
	StoreFactory store2.Factory
}

func NewServer(options Options) (*Server, error) {

	genericServer, err := generic.NewGenericServer(options.ServerOptions, "public")
	if err != nil {
		return nil, err
	}

	container := genericServer.Container

	databaseStore := store2.NewDatabaseStore(options.StoreFactory)
	databaseHandler := database.NewHandler(databaseStore)
	container.Add(databaseHandler.WebService())

	folderStore := store2.NewFolderStore(options.StoreFactory)
	folderHandler := folder.NewHandler(folderStore)
	container.Add(folderHandler.WebService())

	formStore := store2.NewFormStore(options.StoreFactory)
	formHandler := form.NewHandler(formStore)
	container.Add(formHandler.WebService())

	recordStore := store2.NewRecordStore(options.StoreFactory, formStore)
	recordHandler := record.NewHandler(recordStore)
	container.Add(recordHandler.WebService())

	s := &Server{
		Server: genericServer,
	}

	return s, nil
}
