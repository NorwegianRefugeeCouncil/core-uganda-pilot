package cmd

import (
	authnzapiserver "github.com/nrc-no/core/pkg/server/authnzapi"
	"github.com/nrc-no/core/pkg/server/authnzbouncer"
	formsapiserver "github.com/nrc-no/core/pkg/server/formsapi"
	"github.com/nrc-no/core/pkg/server/login"
	client2 "github.com/nrc-no/core/pkg/zanzibar"
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

		client := client2.NewZanzibarClient(coreOptions.ZanzibarClientConfig)

		if err := serveFormsApi(ctx,
			formsapiserver.Options{
				ServerOptions: coreOptions.Serve.FormsApi,
				StoreFactory:  factory,
				ZanzibarClient: client,
			}); err != nil {
			return err
		}
		if err := serveAuthnzApi(ctx,
			authnzapiserver.Options{
				ServerOptions: coreOptions.Serve.AuthnzApi,
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
			authnzbouncer.Options{
				ServerOptions: coreOptions.Serve.AuthnzBouncer,
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
