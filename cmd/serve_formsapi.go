package cmd

import (
	"context"
	formsapiserver "github.com/nrc-no/core/pkg/server/formsapi"
	"github.com/spf13/cobra"
)

// servePublicCmd represents the public command
var servePublicCmd = &cobra.Command{
	Use:   "forms-api",
	Short: "starts the public server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initStoreFactory(); err != nil {
			return err
		}
		if err := servePublic(ctx,
			formsapiserver.Options{
				ServerOptions: coreOptions.Serve.FormsApi,
				StoreFactory:  factory,
			}); err != nil {
			return err
		}
		<-doneSignal
		return nil
	},
}

func init() {
	serveCmd.AddCommand(servePublicCmd)
}

func servePublic(ctx context.Context, options formsapiserver.Options) error {
	server, err := formsapiserver.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
