package cmd

import (
	"github.com/nrc-no/core/cmd/devinit"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devTunnelsCmd = &cobra.Command{
	Use:   "tunnels",
	Short: "starts tunnels",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := devinit.StartTunnels(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	devCmd.AddCommand(devTunnelsCmd)
}
