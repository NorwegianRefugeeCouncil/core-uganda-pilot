package cmd

import (
	"github.com/nrc-no/core/pkg/seeder"
	"github.com/spf13/cobra"
)

var seedColombia = &cobra.Command{
	Use:   "colombia",
	Short: "Seed the database with default forms and fields for the Colombia context",
	RunE:  seedRun(seeder.ColombiaContext),
}

func init() {
	seedCmd.AddCommand(seedColombia)
}
