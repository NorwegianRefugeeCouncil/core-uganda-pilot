package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/server/auth"
	"github.com/spf13/cobra"
)

// serveAuthnzBouncerCmd represents the admin command
var serveAuthnzBouncerCmd = &cobra.Command{
	Use:   "bouncer",
	Short: "starts the authnz-bouncer server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := serveAuthnzBouncer(ctx,
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
	serveCmd.AddCommand(serveAuthnzBouncerCmd)
}

func serveAuthnzBouncer(ctx context.Context, options auth.Options) error {
	server, err := auth.NewServer(ctx, options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
