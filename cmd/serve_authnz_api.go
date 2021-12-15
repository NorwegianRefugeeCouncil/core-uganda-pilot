package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/server/authnzapi"
	"github.com/spf13/cobra"
)

// serveAdminCmd represents the admin command
var serveAdminCmd = &cobra.Command{
	Use:   "authnz-api",
	Short: "starts the authnz-api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initStoreFactory(); err != nil {
			return err
		}
		if err := serveAuthnzApi(ctx,
			authnzapi.Options{
				ServerOptions: coreOptions.Serve.AuthnzApi,
				StoreFactory:  factory,
				HydraAdmin:    coreOptions.Hydra.Admin.AdminClient(),
			}); err != nil {
			return err
		}
		<-doneSignal
		return nil
	},
}

func init() {
	serveCmd.AddCommand(serveAdminCmd)
}

func serveAuthnzApi(ctx context.Context, options authnzapi.Options) error {
	server, err := authnzapi.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
