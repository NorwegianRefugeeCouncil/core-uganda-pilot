package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/server/admin"
	"github.com/spf13/cobra"
)

// serveAdminCmd represents the admin command
var serveAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "starts the admin server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initStoreFactory(); err != nil {
			return err
		}
		if err := serveAdmin(ctx,
			admin.Options{
				ServerOptions: coreOptions.Serve.Admin,
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

func serveAdmin(ctx context.Context, options admin.Options) error {
	server, err := admin.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
