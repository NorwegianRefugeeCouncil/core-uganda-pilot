package cmd

import (
	"github.com/nrc-no/core/pkg/server/admin"
	"github.com/nrc-no/core/pkg/server/public"
	"github.com/spf13/cobra"
)

// serveAllCmd represents the core serve all command
var serveAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Starts the admin, public and login servers",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := servePublic(serveCtx,
			public.Options{
				ServerOptions: coreOptions.Serve.Public,
				StoreFactory:  factory,
			}); err != nil {
			return err
		}

		if err := serveAdmin(serveCtx,
			admin.Options{
				ServerOptions: coreOptions.Serve.Admin,
				StoreFactory:  factory,
			}); err != nil {
			return err
		}

		<-doneSignal
		return nil
	},
}

func init() {
	serveCmd.AddCommand(serveAllCmd)
}
