package cmd

import (
	"context"

	coreDBServer "github.com/nrc-no/core/pkg/server/core-db"
	"github.com/spf13/cobra"
)

var serveCoreDBCmd = &cobra.Command{
	Use:   "core-db-api",
	Short: "starts the core-db-api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initStoreFactory(); err != nil {
			return err
		}
		if err := serveCoreDBApi(ctx,
			coreDBServer.Options{
				ServerOptions: coreOptions.Serve.CoreDB,
				StoreFactory:  factory,
			}); err != nil {
			return err
		}
		<-doneSignal
		return nil
	},
}

func init() {
	serveCmd.AddCommand(serveCoreDBCmd)
}

func serveCoreDBApi(ctx context.Context, options coreDBServer.Options) error {
	server, err := coreDBServer.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
