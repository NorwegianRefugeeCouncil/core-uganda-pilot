package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/server/auth"
	"github.com/spf13/cobra"
)

// serveAuthCmd represents the admin command
var serveAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "starts the auth server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := serveAuth(ctx,
			auth.Options{
				ServerOptions: coreOptions.Serve.Auth,
				HydraPublic:   coreOptions.Hydra.Public.PublicClient(),
				HydraAdmin:    coreOptions.Hydra.Admin.AdminClient(),
			}); err != nil {
			return err
		}
		<-doneSignal
		return nil
	},
}

func init() {
	serveCmd.AddCommand(serveAuthCmd)
}

func serveAuth(ctx context.Context, options auth.Options) error {
	server, err := auth.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
