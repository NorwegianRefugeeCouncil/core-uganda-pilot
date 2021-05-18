package cmd

import (
	"github.com/nrc-no/core/apps/api/cmd/apiserver/app/options"
	"github.com/nrc-no/core/apps/api/pkg/api/defaultscheme"
	"github.com/nrc-no/core/apps/api/pkg/controlplane"
	"github.com/nrc-no/core/apps/api/pkg/server"
	"github.com/nrc-no/core/apps/api/pkg/server/nrc_apiserver"
	serverstorage "github.com/nrc-no/core/apps/api/pkg/server/storage"
	"github.com/spf13/cobra"
)

func NewAPIServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use:          "apiserver",
		Long:         "Core API Server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			completedServerRunOptions, err := complete(s)
			if err != nil {
				return err
			}

			//if errs := completedServerRunOptions.Validate(); len(errs) != 0 {
			//  return exceptions.NewAggregate(errs)
			//}

			return Run(completedServerRunOptions, server.SetupSignalHandler())

		},
	}
}

func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	server, err := CreateServerChain(completeOptions, stopCh)
	if err != nil {
		return err
	}

	prepared, err := server.PrepareRun()
	if err != nil {
		return err
	}

	return prepared.Run(stopCh)
}

func CreateServerChain(completeOptions completedServerRunOptions, ch <-chan struct{}) (interface{}, interface{}) {

}

// completedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type completedServerRunOptions struct {
	*options.ServerRunOptions
}

func complete(s *options.ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions
	options.ServerRunOptions = s
	return options, nil
}

func CreateApiServerConfig(s completedServerRunOptions) (*controlplane.Config, error) {
	genericConfig, storageFactory, err := buildGenericConfig(s.ServerRunOptions)
	if err != nil {
		return nil, err
	}

	config := &controlplane.Config{
		GenericConfig: genericConfig,
	}

}

func buildGenericConfig(s *options.ServerRunOptions) (
	config *server.Config,
	storageFactory *serverstorage.DefaultStorageFactory,
	lastErr error,
) {

	config = server.NewConfig(defaultscheme.Codecs)
	config.MergedResourceConfig = controlplane.DefaultAPIResourceConfigSource()

	storageFactoryConfig := nrc_apiserver.NewStorageFactoryConfig()
	storageFactoryConfig.APIResourceConfig = config.MergedResourceConfig
	completedStorageFactoryConfig, err := storageFactoryConfig.Complete(s.Mongo)
	if err != nil {
		lastErr = err
		return
	}

	storageFactory, lastErr = completedStorageFactoryConfig.New()
	if lastErr != nil {
		return
	}

	if lastErr = s.Mongo.ApplyWithStorageFactoryTo(storageFactory, config); lastErr != nil {
		return
	}

	return

}
