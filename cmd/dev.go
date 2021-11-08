package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "developer tools",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("admin called")
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
