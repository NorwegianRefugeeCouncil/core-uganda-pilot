package app

import (
	"context"
	"github.com/nrc-no/core/api/pkg/server/options"
	"github.com/spf13/cobra"
	"k8s.io/apiserver/pkg/server"
)

func NewStartCoreServer(defaults *options.Options) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch core server",
		Long:  "Launch core server",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			stopCh := server.SetupSignalHandler()
			go func() {
				<-stopCh
				cancel()
			}()

			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(ctx); err != nil {
				return err
			}
			return nil
		},
	}
	o.AddFlags(cmd.Flags())
	return cmd
}
