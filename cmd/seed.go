package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	endpoint string
	seedCmd  = &cobra.Command{
		Use: "seed",
	}
)

func init() {
	seedCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "the endpoint of the Core API, example http://localhost:9000")
	rootCmd.AddCommand(seedCmd)
}
