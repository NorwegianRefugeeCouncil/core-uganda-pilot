package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	hostUri string

	seedCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with folders and forms",
		Long:  `Seed the database with folders and forms`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return storeSetup()
		},
	}
)

func init() {
	seedCmd.PersistentFlags().StringVarP(&hostUri, "host-uri", "h", "", "the URI of the Core API, example http://localhost:9000")
	rootCmd.AddCommand(seedCmd)
}
