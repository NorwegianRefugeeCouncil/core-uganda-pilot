package cmd

import (
	"github.com/nrc-no/core/pkg/server/admin"
	"github.com/nrc-no/core/pkg/server/auth"
	"github.com/nrc-no/core/pkg/server/forms"
	"github.com/nrc-no/core/pkg/server/login"
	"github.com/spf13/cobra"
)

// serveAllCmd represents the core serve all command
var serveAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Starts the admin, public and login servers",
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
		if err := serveAdmin(ctx,
			admin.Options{
				ServerOptions: coreOptions.Serve.Admin,
				StoreFactory:  factory,
				HydraAdmin:    coreOptions.Hydra.Admin.AdminClient(),
			}); err != nil {
			return err
		}
		if err := serveLogin(ctx,
			login.Options{
				ServerOptions: coreOptions.Serve.Login,
				StoreFactory:  factory,
				HydraAdmin:    coreOptions.Hydra.Admin.AdminClient(),
			}); err != nil {
			return err
		}
		if err := serveAuthnzBouncer(ctx,
			auth.Options{
				ServerOptions: coreOptions.Serve.Auth,
				HydraAdmin:    coreOptions.Hydra.Admin.AdminClient(),
				HydraPublic:   coreOptions.Hydra.Public.PublicClient(),
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
