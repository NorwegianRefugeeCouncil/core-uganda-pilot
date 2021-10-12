package main

import (
	"context"
	"flag"
	"github.com/nrc-no/core/pkg/server"
	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()
	options := server.NewOptions()
	cmd := launchCommand(ctx, options)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
	<-ctx.Done()
}

func launchCommand(ctx context.Context, defaults *server.Options) *cobra.Command {

	o := *defaults

	cmd := &cobra.Command{
		Use:   "core-server",
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
