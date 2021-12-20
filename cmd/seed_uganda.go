package cmd

import (
	"github.com/spf13/cobra"
)

var seedUganda = &cobra.Command{
	Use:   "uganda",
	Short: "Seed the database with default forms and fields for the Uganda context",
	RunE: func(cmd *cobra.Command, args []string) error {
		// return seed.Seed(seeder.UgandaContext)
		return nil
	},
}

func init() {
	seedCmd.AddCommand(seedUganda)
}
