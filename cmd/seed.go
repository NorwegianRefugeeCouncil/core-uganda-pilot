package cmd

import (
	"github.com/nrc-no/core/pkg/seeder"
	"github.com/spf13/cobra"
)

var seed *seeder.Seed

var (
	// Used for flags.
	hostUri string
	seedCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with folders and forms",
		Long:  `Seed the database with folders and forms`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
		
			if err = storeSetup(); err != nil {
				return err
			}

			seed, err = seeder.NewSeed(ctx, factory)
			return err
		},
	}
)

func init() {
	seedCmd.PersistentFlags().StringVarP(&hostUri, "host-uri", "h", "", "the URI of the Core API, example http://localhost:9000")
	rootCmd.AddCommand(seedCmd)
}
