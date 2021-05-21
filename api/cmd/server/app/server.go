package app

import (
	"context"
	"github.com/nrc-no/core/api/pkg/server/options"
	"github.com/spf13/cobra"
)

func NewStartCoreServer(defaults *options.Options, ctx context.Context) *cobra.Command {
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
			if err := o.Run(ctx); err != nil {
				return err
			}
			return nil
		},
	}
	o.AddFlags(cmd.Flags())
	return cmd
}
