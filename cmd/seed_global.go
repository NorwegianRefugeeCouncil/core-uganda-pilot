package cmd

import (
	"github.com/nrc-no/core/pkg/seeder"
	"github.com/spf13/cobra"
)

var seedGlobal = &cobra.Command{
	Use:   "global",
	Short: "Seed the database with default forms and fields for the Global context",
	RunE: func(cmd *cobra.Command, args []string) error {
		return seed.Seed(seeder.GlobalContext)
	},
}

func init() {
	seedCmd.AddCommand(seedGlobal)
}
