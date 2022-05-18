package cmd

import (
	"context"

	"github.com/nrc-no/core/pkg/server/data"
	"github.com/spf13/cobra"
)

// serveDataCmd represents the data command
var serveDataCmd = &cobra.Command{
	Use:   "data",
	Short: "starts the data server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := serveDb(ctx,
			data.Options{
				ServerOptions: coreOptions.Serve.Data,
			}); err != nil {
			return err
		}
		<-doneSignal
		return nil
	},
}

func init() {
	serveCmd.AddCommand(serveDataCmd)
}

func serveDb(ctx context.Context, options data.Options) error {
	server, err := data.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
