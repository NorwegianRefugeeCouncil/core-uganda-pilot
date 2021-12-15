package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/server/forms"
	"github.com/spf13/cobra"
)

// servePublicCmd represents the public command
var servePublicCmd = &cobra.Command{
	Use:   "forms",
	Short: "starts the public server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initStoreFactory(); err != nil {
			return err
		}
		if err := servePublic(ctx,
			forms.Options{
				ServerOptions: coreOptions.Serve.Public,
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

func servePublic(ctx context.Context, options forms.Options) error {
	server, err := forms.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
