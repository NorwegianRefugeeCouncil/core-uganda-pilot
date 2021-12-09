package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/server/login"
	"github.com/spf13/cobra"
)

// serveLoginCmd represents the login command
var serveLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "starts the login server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := serveLogin(ctx,
			login.Options{
				ServerOptions: coreOptions.Serve.Login,
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
	serveCmd.AddCommand(serveLoginCmd)
}

func serveLogin(ctx context.Context, options login.Options) error {
	server, err := login.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
