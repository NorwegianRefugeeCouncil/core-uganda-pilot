package cmd

import (
	"github.com/nrc-no/core/cmd/devinit"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := devinit.Init(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	devCmd.AddCommand(devInitCmd)
}
