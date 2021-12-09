package cmd

import (
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "developer tools",
}

func init() {
	rootCmd.AddCommand(devCmd)
}
