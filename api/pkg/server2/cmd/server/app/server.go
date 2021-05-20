package app

import (
	"github.com/nrc-no/core/api/pkg/server2/options"
	"github.com/spf13/cobra"
)

func NewStartCoreServer(defaults *options.Options, stopCh <-chan struct{}) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch core server",
		Long:  "Launch core server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(stopCh); err != nil {
				return err
			}
			return nil
		},
	}
	o.AddFlags(cmd.Flags())
	return cmd
}
