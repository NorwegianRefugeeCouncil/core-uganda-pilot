package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/server/authnzbouncer"
	"github.com/spf13/cobra"
)

// serveAuthnzBouncerCmd represents the admin command
var serveAuthnzBouncerCmd = &cobra.Command{
	Use:   "authnz-bouncer",
	Short: "starts the authnz-bouncer server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := serveAuthnzBouncer(ctx,
			authnzbouncer.Options{
				ServerOptions: coreOptions.Serve.AuthnzBouncer,
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

func serveAuthnzBouncer(ctx context.Context, options authnzbouncer.Options) error {
	server, err := authnzbouncer.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
