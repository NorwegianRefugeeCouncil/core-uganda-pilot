package app

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/nrc-no/core-kafka/pkg/server"
	"github.com/spf13/cobra"
	"log"
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

			// load .env file
			err = godotenv.Load(".env")
			if err != nil {
				log.Fatalf("Error loading .env file")
			}

			_ = completedOptions.New(ctx)

			return nil
		},
	}

	o.Flags(cmd.Flags())
	return cmd

}
