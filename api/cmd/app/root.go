package app

import (
	"context"
	"github.com/nrc-no/core/pkg/server"
	"github.com/spf13/cobra"
)

func LaunchCommand(ctx context.Context, defaults *server.Options) *cobra.Command {

	o := *defaults

	cmd := &cobra.Command{
		Use:   "core-kafka",
		Short: "Core Server",
		Long:  `Core Server`,
		RunE: func(cmd *cobra.Command, args []string) error {

			completedOptions, err := o.Complete(ctx)
			if err != nil {
				panic(err)
			}

			_ = completedOptions.New(ctx)

			return nil
		},
	}

	o.Flags(cmd.Flags())
	return cmd

}
